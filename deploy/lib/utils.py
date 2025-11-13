"""工具函数"""
import os
import re
import socket
import subprocess
from pathlib import Path
from typing import Optional, Dict, Any
from rich.console import Console

console = Console()


def validate_project_name(name: str) -> bool:
    """验证项目名称格式 (a-z, 3-20字符)"""
    pattern = r'^[a-z][a-z0-9-]{2,19}$'
    return bool(re.match(pattern, name))


def check_port_available(port: int) -> bool:
    """检查端口是否可用"""
    try:
        sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        result = sock.connect_ex(('127.0.0.1', port))
        sock.close()
        return result != 0
    except Exception:
        return False


def run_command(cmd: list, check: bool = True) -> bool:
    """运行系统命令"""
    try:
        subprocess.run(cmd, check=check, capture_output=True)
        return True
    except subprocess.CalledProcessError:
        return False


def generate_secret(length: int = 32) -> str:
    """生成随机密钥"""
    import secrets
    import string
    alphabet = string.ascii_letters + string.digits + "!@#$%^&*"
    return ''.join(secrets.choice(alphabet) for _ in range(length))


def get_docker_networks() -> list:
    """获取现有的 Docker 网络列表"""
    try:
        result = subprocess.run(
            ['docker', 'network', 'ls', '--format', '{{.Name}}'],
            capture_output=True,
            text=True,
            check=True
        )
        return result.stdout.strip().split('\n')
    except Exception:
        return []


def docker_network_exists(network_name: str) -> bool:
    """检查 Docker 网络是否存在"""
    networks = get_docker_networks()
    return network_name in networks


def create_docker_network(network_name: str) -> bool:
    """创建 Docker 网络"""
    try:
        subprocess.run(
            ['docker', 'network', 'create', network_name],
            check=True,
            capture_output=True
        )
        return True
    except Exception:
        return False


def test_database_connection(db_type: str, host: str, port: int, 
                            user: str, password: str, db_name: str) -> bool:
    """测试数据库连接"""
    try:
        if db_type == 'postgres':
            import psycopg2
            conn = psycopg2.connect(
                host=host,
                port=port,
                user=user,
                password=password,
                database=db_name
            )
            conn.close()
            return True
        elif db_type == 'mysql':
            import mysql.connector
            conn = mysql.connector.connect(
                host=host,
                port=port,
                user=user,
                password=password,
                database=db_name
            )
            conn.close()
            return True
        elif db_type == 'sqlite3':
            import sqlite3
            conn = sqlite3.connect(db_name)
            conn.close()
            return True
    except Exception as e:
        console.print(f"[red]数据库连接失败: {e}[/red]")
        return False


def test_redis_connection(host: str, port: int, password: str, db: int) -> bool:
    """测试 Redis 连接"""
    try:
        import redis
        r = redis.Redis(
            host=host,
            port=port,
            password=password if password else None,
            db=db,
            decode_responses=True
        )
        r.ping()
        return True
    except Exception as e:
        console.print(f"[red]Redis 连接失败: {e}[/red]")
        return False


def docker_compose_exists() -> bool:
    """检查 docker-compose 是否已安装"""
    try:
        subprocess.run(
            ['docker-compose', '--version'],
            check=True,
            capture_output=True
        )
        return True
    except Exception:
        return False


def pull_docker_image(image: str) -> bool:
    """拉取 Docker 镜像"""
    try:
        console.print(f"[cyan]正在拉取镜像: {image}[/cyan]")
        subprocess.run(
            ['docker', 'pull', image],
            check=True,
            capture_output=False
        )
        return True
    except Exception as e:
        console.print(f"[red]镜像拉取失败: {e}[/red]")
        return False


def docker_image_exists(image: str) -> bool:
    """检查 Docker 镜像是否存在"""
    try:
        subprocess.run(
            ['docker', 'image', 'inspect', image],
            check=True,
            capture_output=True
        )
        return True
    except Exception:
        return False


def generate_docker_compose(config: dict, output_path: str = "docker-compose.yml") -> bool:
    """生成 docker-compose.yml 文件"""
    try:
        # 读取模板文件
        template_path = os.path.join(os.path.dirname(__file__), "..", "templates", "docker-compose.tpl")

        if not os.path.exists(template_path):
            console.print(f"[red]✗ 模板文件不存在: {template_path}[/red]")
            return False

        with open(template_path, 'r', encoding='utf-8') as f:
            template = f.read()

        # 替换模板中的占位符
        content = template
        for key, value in config.items():
            placeholder = "{" + key + "}"
            content = content.replace(placeholder, str(value))

        # 写入文件
        with open(output_path, 'w', encoding='utf-8') as f:
            f.write(content)

        console.print(f"[green]✓ docker-compose.yml 已生成: {output_path}[/green]")
        return True
    except Exception as e:
        console.print(f"[red]✗ 生成 docker-compose.yml 失败: {e}[/red]")
        return False

