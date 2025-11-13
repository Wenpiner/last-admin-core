.PHONY: help install test demo clean venv

help:
	@echo "Last Admin 安装脚本 - 可用命令"
	@echo ""
	@echo "  make install    - 运行安装脚本"
	@echo "  make test       - 运行单元测试"
	@echo "  make demo       - 运行演示脚本"
	@echo "  make venv       - 创建虚拟环境"
	@echo "  make clean      - 清理临时文件"
	@echo "  make help       - 显示此帮助信息"

venv:
	python3 -m venv venv
	. venv/bin/activate && pip install -q -r requirements.txt
	@echo "✓ 虚拟环境已创建"

install: venv
	. venv/bin/activate && python3 install.py

test: venv
	. venv/bin/activate && python3 test_install.py

demo: venv
	. venv/bin/activate && python3 demo_install.py

clean:
	rm -rf venv __pycache__ lib/__pycache__ .pytest_cache .coverage
	rm -f .env.test *.pyc
	@echo "✓ 临时文件已清理"

