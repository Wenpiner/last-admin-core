#!/usr/bin/env python3
"""测试安装脚本"""

import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).parent))

from lib.config import ConfigManager
from lib.utils import (
    validate_project_name, check_port_available, generate_secret,
    get_docker_networks
)


def test_config_manager():
    """测试配置管理器"""
    print("测试 ConfigManager...")
    config = ConfigManager(".env.test")
    config.create_if_not_exists()
    
    config.set("TEST_KEY", "test_value")
    assert config.get("TEST_KEY") == "test_value", "配置保存失败"
    
    print("✓ ConfigManager 测试通过")


def test_validators():
    """测试验证器"""
    print("测试验证器...")
    
    # 项目名称验证
    assert validate_project_name("myproject") == True
    assert validate_project_name("my-project") == True
    assert validate_project_name("MyProject") == False  # 大写不允许
    assert validate_project_name("my") == False  # 太短
    assert validate_project_name("a" * 21) == False  # 太长
    
    print("✓ 验证器测试通过")


def test_utils():
    """测试工具函数"""
    print("测试工具函数...")
    
    # 生成密钥
    secret = generate_secret(32)
    assert len(secret) == 32, "密钥长度不正确"
    
    # 检查端口
    available = check_port_available(9999)
    print(f"  端口 9999 可用: {available}")
    
    # 获取 Docker 网络
    networks = get_docker_networks()
    print(f"  Docker 网络数: {len(networks)}")
    
    print("✓ 工具函数测试通过")


def main():
    """运行所有测试"""
    print("=" * 50)
    print("Last Admin 安装脚本测试")
    print("=" * 50 + "\n")
    
    try:
        test_config_manager()
        test_validators()
        test_utils()
        
        print("\n" + "=" * 50)
        print("✓ 所有测试通过")
        print("=" * 50)
        
        # 清理测试文件
        Path(".env.test").unlink(missing_ok=True)
        
    except Exception as e:
        print(f"\n✗ 测试失败: {e}")
        sys.exit(1)


if __name__ == "__main__":
    main()

