#!/usr/bin/env python3
"""
Last Admin 安装脚本
支持 Linux、macOS、Windows
"""

import sys
import os
from pathlib import Path
from rich.console import Console
from rich.table import Table
from rich.panel import Panel

# 添加 lib 目录到 Python 路径
sys.path.insert(0, str(Path(__file__).parent))

from lib.config import ConfigManager
from lib.prompts import PromptManager
from lib.utils import (
    generate_secret, get_docker_networks, docker_network_exists,
    create_docker_network, validate_project_name, docker_compose_exists,
    pull_docker_image, docker_image_exists, generate_docker_compose
)

console = Console()


class Installer:
    """安装器主类"""

    def __init__(self, install_dir: str = None):
        # 设置安装目录
        if install_dir:
            self.install_dir = install_dir
        else:
            self.install_dir = os.getcwd()

        # 确保安装目录存在
        os.makedirs(self.install_dir, exist_ok=True)

        # 配置文件路径
        env_file = os.path.join(self.install_dir, ".env")

        self.config = ConfigManager(env_file)
        self.config.create_if_not_exists()
        self.prompts = PromptManager(self.config)
        self.data = {}
    
    def print_header(self):
        """打印标题"""
        console.print(Panel(
            "[bold green]Last Admin 安装向导[/bold green]",
            expand=False
        ))
    
    def stage_1_project_info(self):
        """第 1 阶段：项目基本信息"""
        console.print("\n[bold cyan]第 1 阶段：项目基本信息[/bold cyan]")
        
        project_name = self.prompts.prompt_project_name()
        self.config.set("PROJECT_NAME", project_name)
        self.data["PROJECT_NAME"] = project_name
        
        deploy_env = self.prompts.prompt_deploy_env()
        self.config.set("DEPLOY_ENV", deploy_env)
        self.data["DEPLOY_ENV"] = deploy_env
        
        console.print("[green]✓ 项目基本信息已保存[/green]")
    
    def stage_2_docker_network(self):
        """第 2 阶段：Docker 网络配置"""
        console.print("\n[bold cyan]第 2 阶段：Docker 网络配置[/bold cyan]")
        
        existing_networks = get_docker_networks()
        network_choice = self.prompts.prompt_docker_network(existing_networks)
        
        if network_choice == "创建新网络":
            network_name = f"{self.data['PROJECT_NAME']}-network"
            if not docker_network_exists(network_name):
                if create_docker_network(network_name):
                    console.print(f"[green]✓ Docker 网络已创建: {network_name}[/green]")
                else:
                    console.print(f"[red]✗ Docker 网络创建失败[/red]")
                    return False
            self.data["DOCKER_NETWORK"] = network_name
        else:
            self.data["DOCKER_NETWORK"] = network_choice
        
        self.config.set("DOCKER_NETWORK", self.data["DOCKER_NETWORK"])
        return True
    
    def stage_3_docker_images(self):
        """第 3 阶段：Docker 镜像配置"""
        console.print("\n[bold cyan]第 3 阶段：Docker 镜像配置[/bold cyan]")
        
        api_repo, api_tag = self.prompts.prompt_api_image()
        self.config.set("API_IMAGE_REPO", api_repo)
        self.config.set("API_IMAGE_TAG", api_tag)
        self.data["API_IMAGE_REPO"] = api_repo
        self.data["API_IMAGE_TAG"] = api_tag
        
        rpc_repo, rpc_tag = self.prompts.prompt_rpc_image()
        self.config.set("RPC_IMAGE_REPO", rpc_repo)
        self.config.set("RPC_IMAGE_TAG", rpc_tag)
        self.data["RPC_IMAGE_REPO"] = rpc_repo
        self.data["RPC_IMAGE_TAG"] = rpc_tag
        
        console.print("[green]✓ Docker 镜像配置已保存[/green]")
    
    def stage_4_deploy_mode(self):
        """第 4 阶段：组件部署方案"""
        console.print("\n[bold cyan]第 4 阶段：组件部署方案[/bold cyan]")
        
        db_mode = self.prompts.prompt_deploy_mode("database")
        self.config.set("DB_DEPLOY_MODE", db_mode)
        self.data["DB_DEPLOY_MODE"] = db_mode
        
        redis_mode = self.prompts.prompt_deploy_mode("redis")
        self.config.set("REDIS_DEPLOY_MODE", redis_mode)
        self.data["REDIS_DEPLOY_MODE"] = redis_mode
        
        console.print("[green]✓ 组件部署方案已保存[/green]")
    
    def stage_5_ports(self):
        """第 5 阶段：端口配置"""
        console.print("\n[bold cyan]第 5 阶段：端口配置[/bold cyan]")
        
        api_port = self.prompts.prompt_port(
            "API 服务端口",
            default=int(self.config.get("API_PORT", "8889")),
            check_available=True
        )
        self.config.set("API_PORT", str(api_port))
        self.data["API_PORT"] = str(api_port)
        
        rpc_port = self.prompts.prompt_port(
            "RPC 服务端口",
            default=int(self.config.get("RPC_PORT", "8080")),
            check_available=True
        )
        self.config.set("RPC_PORT", str(rpc_port))
        self.data["RPC_PORT"] = str(rpc_port)
        
        if self.data["DB_DEPLOY_MODE"] == "docker":
            db_port = self.prompts.prompt_port(
                "数据库端口",
                default=int(self.config.get("DB_PORT", "5432")),
                check_available=True
            )
        else:
            db_port = self.prompts.prompt_port(
                "数据库端口",
                default=int(self.config.get("DB_PORT", "5432")),
                check_available=False
            )
        self.config.set("DB_PORT", str(db_port))
        self.data["DB_PORT"] = str(db_port)
        
        if self.data["REDIS_DEPLOY_MODE"] == "docker":
            redis_port = self.prompts.prompt_port(
                "Redis 端口",
                default=int(self.config.get("REDIS_PORT", "6379")),
                check_available=True
            )
        else:
            redis_port = self.prompts.prompt_port(
                "Redis 端口",
                default=int(self.config.get("REDIS_PORT", "6379")),
                check_available=False
            )
        self.config.set("REDIS_PORT", str(redis_port))
        self.data["REDIS_PORT"] = str(redis_port)
        
        console.print("[green]✓ 端口配置已保存[/green]")
    
    def stage_6_database(self):
        """第 6 阶段：数据库配置"""
        console.print("\n[bold cyan]第 6 阶段：数据库配置[/bold cyan]")
        
        db_type = self.prompts.prompt_db_type()
        self.config.set("DB_TYPE", db_type)
        self.data["DB_TYPE"] = db_type
        
        db_user = self.prompts.prompt_text(
            "数据库用户",
            default=self.config.get("DB_USER", "postgres")
        )
        self.config.set("DB_USER", db_user)
        self.data["DB_USER"] = db_user
        
        db_password = self.prompts.prompt_text(
            "数据库密码",
            default=self.config.get("DB_PASSWORD", "postgres123")
        )
        self.config.set("DB_PASSWORD", db_password)
        self.data["DB_PASSWORD"] = db_password
        
        db_name = self.prompts.prompt_text(
            "数据库名称",
            default=self.config.get("DB_NAME", f"{self.data['PROJECT_NAME']}_db")
        )
        self.config.set("DB_NAME", db_name)
        self.data["DB_NAME"] = db_name
        
        db_ssl_mode = self.prompts.prompt_choice(
            "SSL 模式",
            ["disable", "require", "prefer"],
            default=self.config.get("DB_SSL_MODE", "disable")
        )
        self.config.set("DB_SSL_MODE", db_ssl_mode)
        self.data["DB_SSL_MODE"] = db_ssl_mode
        
        if self.data["DB_DEPLOY_MODE"] == "docker":
            # Docker 部署时，使用 Docker 服务名称作为 host
            db_host = f"postgres-{self.data['PROJECT_NAME']}"
            console.print(f"[cyan]ℹ Docker 部署模式，数据库主机已自动设置为: {db_host}[/cyan]")
        else:
            # 外部部署时，需要用户输入数据库主机
            db_host = self.prompts.prompt_text(
                "数据库主机 (外部部署)",
                default=self.config.get("DB_HOST", "localhost")
            )
        self.config.set("DB_HOST", db_host)
        self.data["DB_HOST"] = db_host
        
        console.print("[green]✓ 数据库配置已保存[/green]")
    
    def stage_7_redis(self):
        """第 7 阶段：Redis 配置"""
        console.print("\n[bold cyan]第 7 阶段：Redis 配置[/bold cyan]")
        
        redis_password = self.prompts.prompt_text(
            "Redis 密码",
            default=self.config.get("REDIS_PASSWORD", "redis123")
        )
        self.config.set("REDIS_PASSWORD", redis_password)
        self.data["REDIS_PASSWORD"] = redis_password
        
        redis_db = self.prompts.prompt_text(
            "Redis 数据库编号",
            default=self.config.get("REDIS_DB", "0")
        )
        self.config.set("REDIS_DB", redis_db)
        self.data["REDIS_DB"] = redis_db
        
        redis_pool_size = self.prompts.prompt_text(
            "Redis 连接池大小",
            default=self.config.get("REDIS_POOL_SIZE", "10")
        )
        self.config.set("REDIS_POOL_SIZE", redis_pool_size)
        self.data["REDIS_POOL_SIZE"] = redis_pool_size
        
        if self.data["REDIS_DEPLOY_MODE"] == "docker":
            # Docker 部署时，使用 Docker 服务名称作为 host，内部端口固定为 6379
            redis_host = f"redis-{self.data['PROJECT_NAME']}:6379"
            console.print(f"[cyan]ℹ Docker 部署模式，Redis 主机已自动设置为: {redis_host}[/cyan]")
        else:
            # 外部部署时，需要用户输入 Redis 主机
            redis_host = self.prompts.prompt_text(
                "Redis 主机 (外部部署，格式: host:port)",
                default=self.config.get("REDIS_HOST", "localhost:6379")
            )
        self.config.set("REDIS_HOST", redis_host)
        self.data["REDIS_HOST"] = redis_host
        
        console.print("[green]✓ Redis 配置已保存[/green]")
    
    def stage_8_auth(self):
        """第 8 阶段：认证和密钥配置"""
        console.print("\n[bold cyan]第 8 阶段：认证和密钥配置[/bold cyan]")
        
        auth_secret = self.config.get("AUTH_ACCESS_SECRET", "")
        if not auth_secret:
            auth_secret = generate_secret()
            console.print(f"[yellow]已自动生成 JWT 密钥[/yellow]")
        self.config.set("AUTH_ACCESS_SECRET", auth_secret)
        self.data["AUTH_ACCESS_SECRET"] = auth_secret
        
        auth_expire = self.prompts.prompt_text(
            "Token 过期时间 (秒)",
            default=self.config.get("AUTH_ACCESS_EXPIRE", "360000")
        )
        self.config.set("AUTH_ACCESS_EXPIRE", auth_expire)
        self.data["AUTH_ACCESS_EXPIRE"] = auth_expire
        
        oauth_secret = self.config.get("OAUTH_STATE_SECRET", "")
        if not oauth_secret:
            oauth_secret = generate_secret()
            console.print(f"[yellow]已自动生成 OAuth 密钥[/yellow]")
        self.config.set("OAUTH_STATE_SECRET", oauth_secret)
        self.data["OAUTH_STATE_SECRET"] = oauth_secret
        
        console.print("[green]✓ 认证配置已保存[/green]")
    
    def stage_9_captcha(self):
        """第 9 阶段：验证码配置"""
        console.print("\n[bold cyan]第 9 阶段：验证码配置[/bold cyan]")

        captcha_type = self.prompts.prompt_captcha_type()
        self.config.set("CAPTCHA_TYPE", captcha_type)
        self.data["CAPTCHA_TYPE"] = captcha_type

        captcha_store = self.prompts.prompt_choice(
            "验证码存储类型",
            ["memory", "redis"],
            default=self.config.get("CAPTCHA_STORE_TYPE", "redis")
        )
        self.config.set("CAPTCHA_STORE_TYPE", captcha_store)
        self.data["CAPTCHA_STORE_TYPE"] = captcha_store

        console.print("[green]✓ 验证码配置已保存[/green]")
    
    def show_summary(self):
        """显示配置摘要"""
        console.print("\n[bold cyan]配置摘要[/bold cyan]")

        table = Table(title="部署配置")
        table.add_column("配置项", style="cyan")
        table.add_column("值", style="magenta")

        for key, value in sorted(self.data.items()):
            if "PASSWORD" in key or "SECRET" in key:
                display_value = "***" + value[-4:] if len(value) > 4 else "***"
            else:
                display_value = str(value)
            table.add_row(key, display_value)

        console.print(table)

    def stage_10_deploy(self):
        """第 10 阶段：部署服务"""
        console.print("\n[bold cyan]第 10 阶段：部署服务[/bold cyan]")

        # 检查 docker compose
        if not docker_compose_exists():
            console.print("[red]✗ 未检测到 docker compose，请先安装[/red]")
            return False

        console.print("[green]✓ docker compose 已安装[/green]")

        # 生成 docker-compose.yml
        console.print("\n[cyan]正在生成 docker-compose.yml...[/cyan]")
        docker_compose_path = os.path.join(self.install_dir, "docker-compose.yml")
        if not generate_docker_compose(self.data, docker_compose_path):
            console.print("[red]✗ 生成 docker-compose.yml 失败[/red]")
            return False
        console.print(f"[green]✓ docker-compose.yml 已生成: {docker_compose_path}[/green]")

        # 拉取 Docker 镜像
        console.print("\n[cyan]正在拉取 Docker 镜像...[/cyan]")

        api_image = f"{self.data['API_IMAGE_REPO']}:{self.data['API_IMAGE_TAG']}"
        rpc_image = f"{self.data['RPC_IMAGE_REPO']}:{self.data['RPC_IMAGE_TAG']}"

        if not docker_image_exists(api_image):
            if not pull_docker_image(api_image):
                console.print(f"[red]✗ 无法拉取 API 镜像: {api_image}[/red]")
                return False
        else:
            console.print(f"[green]✓ API 镜像已存在: {api_image}[/green]")

        if not docker_image_exists(rpc_image):
            if not pull_docker_image(rpc_image):
                console.print(f"[red]✗ 无法拉取 RPC 镜像: {rpc_image}[/red]")
                return False
        else:
            console.print(f"[green]✓ RPC 镜像已存在: {rpc_image}[/green]")

        console.print("[green]✓ Docker 镜像已准备就绪[/green]")
        return True
    
    def run(self):
        """运行安装流程"""
        try:
            self.print_header()

            self.stage_1_project_info()
            self.stage_2_docker_network()
            self.stage_3_docker_images()
            self.stage_4_deploy_mode()
            self.stage_5_ports()
            self.stage_6_database()
            self.stage_7_redis()
            self.stage_8_auth()
            self.stage_9_captcha()

            self.show_summary()

            if self.prompts.prompt_confirm("\n是否继续部署?", default=True):
                env_file = os.path.join(self.install_dir, ".env")
                console.print(f"\n[green]✓ 配置已保存到: {env_file}[/green]")

                # 执行部署步骤
                if self.stage_10_deploy():
                    console.print("\n[green]✓ 部署完成！[/green]")
                    console.print("\n[cyan]后续步骤:[/cyan]")
                    console.print(f"  1. 进入安装目录: cd {self.install_dir}")
                    console.print("  2. 使用 docker compose 启动服务: docker compose up -d")
                    console.print("  3. 检查服务日志: docker compose logs -f")
                    console.print("  4. 访问 API: http://localhost:8889")
                else:
                    console.print("\n[red]✗ 部署失败，请检查错误信息[/red]")
                    sys.exit(1)
            else:
                env_file = os.path.join(self.install_dir, ".env")
                console.print(f"\n[yellow]⚠ 部署已取消，配置已保存到: {env_file}[/yellow]")

        except KeyboardInterrupt:
            console.print("\n[yellow]⚠ 安装已中断，配置已保存，下次运行时将继续[/yellow]")
            sys.exit(0)
        except Exception as e:
            console.print(f"\n[red]✗ 错误: {e}[/red]")
            sys.exit(1)


if __name__ == "__main__":
    # 从命令行参数获取安装目录
    install_dir = None
    if len(sys.argv) > 1:
        install_dir = sys.argv[1]

    installer = Installer(install_dir)
    installer.run()

