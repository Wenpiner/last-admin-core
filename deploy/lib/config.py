"""配置管理"""
from pathlib import Path
from typing import Optional, Dict, Any
from dotenv import load_dotenv, set_key
import os


class ConfigManager:
    """配置管理器"""
    
    def __init__(self, env_file: str = ".env"):
        self.env_file = Path(env_file)
        self.config: Dict[str, Any] = {}
        self.load_config()
    
    def load_config(self):
        """从 .env 文件加载配置"""
        if self.env_file.exists():
            load_dotenv(self.env_file)
            # 读取所有环境变量
            for key in os.environ:
                self.config[key] = os.environ[key]
    
    def get(self, key: str, default: str = "") -> str:
        """获取配置值"""
        return self.config.get(key, default)
    
    def set(self, key: str, value: str):
        """设置配置值并立即保存"""
        self.config[key] = value
        set_key(str(self.env_file), key, value)
    
    def get_all(self) -> Dict[str, Any]:
        """获取所有配置"""
        return self.config.copy()
    
    def update(self, data: Dict[str, str]):
        """批量更新配置"""
        for key, value in data.items():
            self.set(key, value)
    
    def exists(self) -> bool:
        """检查 .env 文件是否存在"""
        return self.env_file.exists()
    
    def create_if_not_exists(self):
        """如果不存在则创建 .env 文件"""
        if not self.env_file.exists():
            self.env_file.touch()

