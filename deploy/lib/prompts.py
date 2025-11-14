"""交互式提示"""
from typing import Optional, List
from rich.console import Console
from rich.prompt import Prompt, Confirm
from lib.utils import validate_project_name, check_port_available
from lib.config import ConfigManager

console = Console()


class PromptManager:
    """提示管理器"""
    
    def __init__(self, config: ConfigManager):
        self.config = config
    
    def prompt_text(self, prompt: str, default: str = "",
                   validator=None) -> str:
        """文本输入提示"""
        while True:
            try:
                value = Prompt.ask(prompt, default=default)
            except EOFError:
                # 非交互式环境，使用默认值
                if default is not None:
                    console.print(f"[yellow]⚠ 非交互式环境，使用默认值: {default if default else '(空)'}[/yellow]")
                    return default
                else:
                    console.print("[red]✗ 非交互式环境且无默认值，请提供输入[/red]")
                    raise

            if validator and not validator(value):
                console.print("[red]✗ 输入格式不正确，请重试[/red]")
                continue
            return value
    
    def prompt_choice(self, prompt: str, choices: List[str],
                     default: str = "") -> str:
        """选择提示"""
        if not default or default not in choices:
            default = choices[0]

        console.print(f"\n{prompt}")
        for i, choice in enumerate(choices, 1):
            console.print(f"  {i}. {choice}")

        while True:
            try:
                default_index = choices.index(default) + 1
                try:
                    choice_input = Prompt.ask("请选择", default=str(default_index))
                except EOFError:
                    # 非交互式环境，使用默认值
                    console.print(f"[yellow]⚠ 非交互式环境，使用默认值: {default}[/yellow]")
                    return default

                index = int(choice_input) - 1
                if 0 <= index < len(choices):
                    return choices[index]
            except (ValueError, IndexError):
                pass
            console.print("[red]✗ 选择无效，请重试[/red]")
    
    def prompt_port(self, prompt: str, default: int = 0,
                   check_available: bool = True) -> int:
        """端口输入提示"""
        while True:
            try:
                port_str = Prompt.ask(prompt, default=str(default))
            except EOFError:
                # 非交互式环境，使用默认值
                if default and default > 0:
                    console.print(f"[yellow]⚠ 非交互式环境，使用默认值: {default}[/yellow]")
                    return default
                else:
                    console.print("[red]✗ 非交互式环境且无有效默认值，请提供输入[/red]")
                    raise

            try:
                port = int(port_str)
                if port < 1 or port > 65535:
                    console.print("[red]✗ 端口号必须在 1-65535 之间[/red]")
                    continue

                if check_available and not check_port_available(port):
                    console.print(f"[red]✗ 端口 {port} 已被占用，请选择其他端口[/red]")
                    continue

                return port
            except ValueError:
                console.print("[red]✗ 请输入有效的端口号[/red]")
    
    def prompt_confirm(self, prompt: str, default: bool = True) -> bool:
        """确认提示"""
        try:
            return Confirm.ask(prompt, default=default)
        except EOFError:
            # 非交互式环境，使用默认值
            console.print(f"[yellow]⚠ 非交互式环境，使用默认值: {'是' if default else '否'}[/yellow]")
            return default
    
    def prompt_project_name(self) -> str:
        """项目名称提示"""
        default = self.config.get("PROJECT_NAME", "lastadmin")
        return self.prompt_text(
            "项目名称 (a-z, 3-20字符)",
            default=default,
            validator=validate_project_name
        )
    
    def prompt_deploy_env(self) -> str:
        """部署环境提示"""
        default = self.config.get("DEPLOY_ENV", "prod")
        return self.prompt_choice(
            "部署环境",
            ["dev", "test", "prod"],
            default=default
        )
    
    def prompt_docker_network(self, existing_networks: List[str]) -> str:
        """Docker 网络选择提示"""
        choices = ["创建新网络", "自定义输入"] + existing_networks

        # 读取缓存的默认值
        cached_network = self.config.get("DOCKER_NETWORK", "")
        default_choice = "创建新网络"

        # 如果缓存值在选项中，使用缓存值作为默认
        if cached_network in choices:
            default_choice = cached_network

        choice = self.prompt_choice(
            "Docker 网络",
            choices,
            default=default_choice
        )

        if choice == "自定义输入":
            return self.prompt_text(
                "请输入 Docker 网络名称",
                default=self.config.get("DOCKER_NETWORK", "")
            )
        return choice
    
    def prompt_api_image(self) -> tuple:
        """API 镜像提示"""
        repo = self.prompt_text(
            "API 镜像仓库",
            default=self.config.get("API_IMAGE_REPO", "wenpiner/last-admin-api")
        )
        tag = self.prompt_text(
            "API 镜像标签",
            default=self.config.get("API_IMAGE_TAG", "latest")
        )
        return repo, tag

    def prompt_rpc_image(self) -> tuple:
        """RPC 镜像提示"""
        repo = self.prompt_text(
            "RPC 镜像仓库",
            default=self.config.get("RPC_IMAGE_REPO", "wenpiner/last-admin-rpc")
        )
        tag = self.prompt_text(
            "RPC 镜像标签",
            default=self.config.get("RPC_IMAGE_TAG", "latest")
        )
        return repo, tag
    
    def prompt_deploy_mode(self, component: str) -> str:
        """部署方式提示"""
        # 将组件名称转换为中文
        component_cn = {
            "database": "数据库",
            "redis": "Redis"
        }.get(component, component)

        default = self.config.get(f"{component.upper()}_DEPLOY_MODE", "docker")
        return self.prompt_choice(
            f"{component_cn}部署方式",
            ["docker", "external"],
            default=default
        )
    
    def prompt_db_type(self) -> str:
        """数据库类型提示"""
        default = self.config.get("DB_TYPE", "postgres")
        return self.prompt_choice(
            "数据库类型",
            ["postgres", "mysql", "sqlite3"],
            default=default
        )

    def prompt_captcha_type(self) -> str:
        """验证码类型提示"""
        captcha_types = [
            "digit (数字)",
            "string (字符串)",
            "math (数学)",
            "chinese (中文)",
            "audio (音频)",
            "random (随机)"
        ]
        default_type = self.config.get("CAPTCHA_TYPE", "random")
        # 找到对应的显示值
        default_display = next(
            (t for t in captcha_types if t.startswith(default_type)),
            captcha_types[-1]
        )

        choice = self.prompt_choice(
            "验证码类型",
            captcha_types,
            default=default_display
        )
        # 提取类型代码（去掉中文说明）
        return choice.split(" ")[0]

