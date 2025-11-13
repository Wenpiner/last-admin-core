#!/usr/bin/env python3
"""演示安装脚本（自动输入）"""

import sys
import os
from pathlib import Path
from unittest.mock import patch
from io import StringIO

sys.path.insert(0, str(Path(__file__).parent))

from install import Installer


def demo_install():
    """演示安装流程"""

    # 模拟用户输入
    inputs = [
        "myproject",                # 项目名称
        "1",                        # 部署环境 (dev)
        "1",                        # Docker 网络 (创建新网络)
        "wenpiner/last-admin-api",  # API 镜像仓库
        "latest",                   # API 镜像标签
        "wenpiner/last-admin-rpc",  # RPC 镜像仓库
        "latest",                   # RPC 镜像标签
        "1",                        # 数据库部署方式 (docker)
        "2",                        # Redis 部署方式 (external)
        "8889",                     # API 端口
        "8080",                     # RPC 端口
        "5433",                     # 数据库端口 (避免冲突)
        "6379",                     # Redis 端口
        "1",                        # 数据库类型 (postgres)
        "postgres",                 # 数据库用户
        "postgres123",              # 数据库密码
        "myproject_db",             # 数据库名称
        "1",                        # SSL 模式 (disable)
        "redis123",                 # Redis 密码
        "0",                        # Redis 数据库编号
        "10",                       # Redis 连接池大小
        "localhost:6379",           # Redis 主机 (外部模式)
        "360000",                   # Token 过期时间
        "1",                        # 验证码类型 (digit - 数字)
        "2",                        # 验证码存储类型 (redis)
        "y",                        # 是否继续部署
    ]
    
    print("=" * 60)
    print("Last Admin 安装脚本演示")
    print("=" * 60)
    print("\n模拟用户输入，自动完成安装配置...\n")
    
    with patch('builtins.input', side_effect=inputs):
        try:
            installer = Installer()
            installer.run()
            
            print("\n" + "=" * 60)
            print("✓ 演示完成")
            print("=" * 60)
            
            # 显示生成的 .env 文件内容
            if Path(".env").exists():
                print("\n生成的 .env 文件内容:")
                print("-" * 60)
                with open(".env", "r") as f:
                    for line in f:
                        if "PASSWORD" in line or "SECRET" in line:
                            key, _ = line.split("=", 1)
                            print(f"{key}=***")
                        else:
                            print(line.rstrip())
                print("-" * 60)
        
        except Exception as e:
            print(f"\n✗ 演示失败: {e}")
            import traceback
            traceback.print_exc()
            sys.exit(1)


if __name__ == "__main__":
    demo_install()

