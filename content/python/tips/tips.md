---
date: 2026-09-22
---
### pyenv setup

```bash
# install pyenv on macos
brew install pyenv

# add path for zsh
echo 'export PYENV_ROOT="$HOME/.pyenv"' >> ~/.zshrc
echo '[[ -d $PYENV_ROOT/bin ]] && export PATH="$PYENV_ROOT/bin:$PATH"' >> ~/.zshrc
echo 'eval "$(pyenv init - zsh)"' >> ~/.zshrc

exec "$SHELL"

# show current python version
pyenv version

# show all installed python version
pyenv versions

# show all remote python version
pyenv install --list

# install specific version for python
pyenv install *.*.*

# setup current directory python version
pyenv local *.*.*

# setup globally python version
pyenv global *.*.*
```

### create new project with pyenv

```bash
# change local directory's python version
pyenv local *.*.*

# check current directory's python version
python -V

# create virtual environment
python -m venv .venv

# activate
source .venv/bin/activate

# deactivate
deactivate
```

### pip: common commands

```bash
# check pip version
pip -V

# upgrade pip
pip install --upgrade pip

# install package
pip install [package]

# upgrade package
pip install -U [package]

# uninstall package
pip uninstall [package]

# list all installed package
pip list

# show package details
pip show [package]
```

### pip: write packages to requirements.txt

```bash
# write requirements
pip freeze > requirements.txt

# install from requirements.txt
pip install -r requirements.txt
```

### dotenv

```bash
# install dotenv
pip install python-dotenv

# load dotenv and get env variable
import os
from dotenv import load_dotenv, find_dotenv

load_dotenv(find_dotenv(), override=True)
os.getenv("ENV_VARIABLE_NAME")
```
