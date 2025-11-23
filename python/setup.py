from setuptools import setup, find_packages

setup(
    name="cloudbridge-sdk",
    version="0.1.0",
    description="Official CloudBridge SDK for Python",
    author="2GC CloudBridge",
    author_email="support@2gc.ru",
    packages=find_packages(),
    install_requires=[
        "requests>=2.25.0",
        "python-jose>=3.3.0",
        "websockets>=10.0",
    ],
    python_requires=">=3.8",
)
