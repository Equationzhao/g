import pandas as pd
import streamlit as st

# Set page configuration
st.set_page_config(
    page_title="g - Advanced ls Command Alternative",
    page_icon="📂",
    layout="wide",
    initial_sidebar_state="expanded"
)

# Custom CSS
st.markdown("""
<style>
    .main-title {
        font-size: 3rem;
        font-weight: bold;
        color: #1E88E5;
        margin-bottom: 1rem;
        text-align: center;
    }
    .section-header {
        font-size: 2rem;
        color: #0D47A1;
        border-bottom: 2px solid #90CAF9;
        padding-bottom: 0.5rem;
        margin-top: 2rem;
    }
    .feature-card {
        background-color: #f8f9fa;
        border-radius: 8px;
        padding: 20px;
        box-shadow: 0 4px 6px rgba(0,0,0,0.1);
        height: 100%;
    }
    .feature-title {
        font-weight: bold;
        color: #1565C0;
        margin-bottom: 10px;
    }
    .code-header {
        font-weight: bold;
        margin-top: 1rem;
    }
    .screenshot {
        border-radius: 8px;
        box-shadow: 0 4px 8px rgba(0,0,0,0.2);
    }
    .badge-container {
        display: flex;
        justify-content: center;
        flex-wrap: wrap;
        gap: 10px;
        margin-bottom: 20px;
    }
    .badge-container img {
        height: 25px;
    }
</style>
""", unsafe_allow_html=True)

# Sidebar
with st.sidebar:
    st.title("Navigation")
    st.markdown("---")
    
    st.markdown("### Content")
    st.markdown("- [Key Features](#key-features)")
    st.markdown("- [Screenshots](#screenshots)")
    st.markdown("- [Usage](#usage)")
    st.markdown("- [Installation Guide](#installation-guide)")
    st.markdown("- [Shell Integration](#shell-integration)")
    st.markdown("- [Project Comparison](#project-comparison)")
    
    st.markdown("---")
    st.markdown("### Links")
    st.markdown("- [GitHub Repository](https://github.com/Equationzhao/g)")
    st.markdown("- [Report Issues](https://github.com/Equationzhao/g/issues)")
    st.markdown("- [Theme Documentation](https://github.com/Equationzhao/g/blob/master/docs/Theme.md)")
    st.markdown("- [Manual](https://github.com/Equationzhao/g/blob/master/docs/man.md)")

# Main title
st.markdown('<div class="main-title">🌈 g - Feature-rich ls Alternative</div>', unsafe_allow_html=True)

# Badges
st.markdown('<div class="badge-container">'
            '<img src="https://img.shields.io/github/stars/Equationzhao/g" alt="Stars">'
            '<img src="https://img.shields.io/github/forks/Equationzhao/g" alt="Forks">'
            '<img src="https://img.shields.io/github/issues/Equationzhao/g" alt="Issues">'
            '<img src="https://img.shields.io/github/license/Equationzhao/g" alt="License">'
            '</div>', unsafe_allow_html=True)

# Introduction
st.markdown("""
**g** is a feature-rich, customizable, and cross-platform `ls` alternative. It provides enhanced visuals with type-specific icons, various layout options, and Git status integration.
""")

# Key Features
st.markdown('<div class="section-header" id="key-features">Key Features</div>', unsafe_allow_html=True)

col1, col2 = st.columns(2)

with col1:
    with st.container():
        st.markdown('<div class="feature-card">', unsafe_allow_html=True)
        st.markdown('<div class="feature-title">🎨 Customizable Display</div>', unsafe_allow_html=True)
        st.markdown("Icons and colors specific to file types, easy to customize")
        st.markdown('</div>', unsafe_allow_html=True)

    st.markdown("")  # Spacing

    with st.container():
        st.markdown('<div class="feature-card">', unsafe_allow_html=True)
        st.markdown('<div class="feature-title">🔀 Multiple Layouts</div>', unsafe_allow_html=True)
        st.markdown("Choose from grid, across, byline, zero, comma, table, json, markdown, and tree layouts")
        st.markdown('</div>', unsafe_allow_html=True)

    st.markdown("")  # Spacing

    with st.container():
        st.markdown('<div class="feature-card">', unsafe_allow_html=True)
        st.markdown('<div class="feature-title">🌐 Git Integration</div>', unsafe_allow_html=True)
        st.markdown("View file git-status/repo-status/repo-branch directly in your listings")
        st.markdown('</div>', unsafe_allow_html=True)

with col2:
    with st.container():
        st.markdown('<div class="feature-card">', unsafe_allow_html=True)
        st.markdown('<div class="feature-title">🔄 Advanced Sorting</div>', unsafe_allow_html=True)
        st.markdown("Highly customizable sorting options like version-sort")
        st.markdown('</div>', unsafe_allow_html=True)

    st.markdown("")  # Spacing

    with st.container():
        st.markdown('<div class="feature-card">', unsafe_allow_html=True)
        st.markdown('<div class="feature-title">💻 Cross-Platform</div>', unsafe_allow_html=True)
        st.markdown("Works seamlessly on Linux, Windows, and MacOS")
        st.markdown('</div>', unsafe_allow_html=True)

    st.markdown("")  # Spacing

    with st.container():
        st.markdown('<div class="feature-card">', unsafe_allow_html=True)
        st.markdown('<div class="feature-title">🔍 Fuzzy Path Matching</div>', unsafe_allow_html=True)
        st.markdown("zoxide and fzf like fuzzy path matching")
        st.markdown('</div>', unsafe_allow_html=True)

# Screenshots
st.markdown('<div class="section-header" id="screenshots">Screenshots</div>', unsafe_allow_html=True)
st.image("https://raw.githubusercontent.com/Equationzhao/g/master/asset/screenshot_3.png", 
         caption="g command interface", 
         use_column_width=True, 
         output_format="PNG")

# Usage
st.markdown('<div class="section-header" id="usage">Usage</div>', unsafe_allow_html=True)

st.markdown('<div class="code-header">Basic usage:</div>', unsafe_allow_html=True)
st.code("g path(s)", language="bash")

st.markdown('<div class="code-header">Display icons and long format:</div>', unsafe_allow_html=True)
st.code("g --icon --long path(s)", language="bash")

st.markdown('<div class="code-header">Tree layout:</div>', unsafe_allow_html=True)
st.code("g --tree --long path(s)", language="bash")

more_options = st.expander("View more options")
with more_options:
    st.markdown("""
    ### Layout Options
    - `--grid`: Grid layout (default)
    - `--tree`: Tree layout
    - `--across`: Across layout
    - `--oneline`: One line per entry
    - `--table`: Table layout
    - `--json`: JSON output
    - `--markdown`: Markdown table

    ### Display Options
    - `--icon`: Show file type icons
    - `--long`, `-l`: Long format with details
    - `--all`, `-a`: Show hidden files
    - `--git`: Show git status
    - `--mime`: Show MIME types
    - `--hyperlink`: Enable file hyperlinks

    For all options, see the [manual](https://github.com/Equationzhao/g/blob/master/docs/man.md)
    """)
    
# Custom Themes
theme = st.expander("Custom Themes", expanded=False)
with theme:
    st.markdown("""
    g supports rich theme customization.
    
    For detailed information, see the [Theme Documentation](https://github.com/Equationzhao/g/blob/master/docs/Theme.md).
    
    Basic usage:
    ```bash
    g --theme mytheme.json  # Use custom theme
    ```
    """)

# Installation Guide
st.markdown('<div class="section-header" id="installation-guide">Installation Guide</div>', unsafe_allow_html=True)

install_tabs = st.tabs(["Package Managers", "Pre-built Binaries", "Build from Source"])

with install_tabs[0]:
    col1, col2 = st.columns(2)
    
    with col1:
        st.subheader("Arch Linux")
        st.code("yay -S g-ls", language="bash")
        
        st.subheader("Homebrew")
        st.code("brew install g-ls", language="bash")
        
        st.subheader("MacPort")
        st.code("sudo port install g-ls", language="bash")
    
    with col2:
        st.subheader("Windows (Scoop)")
        st.code("scoop install https://raw.githubusercontent.com/Equationzhao/g/master/scoop/g.json", language="powershell")
        
        st.subheader("WinGet")
        st.info("In development, see [#119](https://github.com/Equationzhao/g/issues/119)")

with install_tabs[1]:
    col1, col2 = st.columns(2)
    
    with col1:
        st.subheader("Install Script")
        st.code('bash -c "$(curl -fsSLk https://raw.githubusercontent.com/Equationzhao/g/master/script/install.sh)"', language="bash")
        
        st.subheader("Uninstall")
        st.code('curl -fsSLk https://raw.githubusercontent.com/Equationzhao/g/master/script/install.sh | bash /dev/stdin -r', language="bash")
    
    with col2:
        st.subheader("Download Packages")
        st.markdown("""
        1. Download from [GitHub releases page](https://github.com/Equationzhao/g/releases)
        2. Extract and add the executable to your PATH
        
        #### Debian/Ubuntu
        ```bash
        sudo dpkg -i g_$version_$arch.deb
        ```
        """)

with install_tabs[2]:
    st.markdown("""
    ### Requirements
    - Go version >= 1.24
    
    ### Install Latest Version
    ```bash
    go install -ldflags="-s -w" github.com/Equationzhao/g@latest
    ```
    
    ### Build from Source
    ```bash
    git clone github.com/Equationzhao/g
    cd g
    go build -ldflags="-s -w"
    # then add the executable file to your PATH
    ```
    """)

# Recommended Terminals
st.markdown("### Recommended Terminals")

term_cols = st.columns(3)

with term_cols[0]:
    st.markdown("#### macOS")
    st.markdown("- [Iterm2](https://iterm2.com/)")
    st.markdown("- [Warp](https://www.warp.dev)")

with term_cols[1]:
    st.markdown("#### Windows")
    st.markdown("- [Windows Terminal](https://github.com/microsoft/terminal)")

with term_cols[2]:
    st.markdown("#### Cross-platform")
    st.markdown("- [Hyper](https://hyper.is/)")
    st.markdown("- [WezTerm](https://wezfurlong.org/wezterm/)")

# Shell Integration
st.markdown('<div class="section-header" id="shell-integration">Shell Integration</div>', unsafe_allow_html=True)

shell_tabs = st.tabs(["Command Completion", "Shell Aliases"])

with shell_tabs[0]:
    st.info("If you install `g` through brew or the install script, the completion is usually installed already.")
    
    completion_tabs = st.tabs(["zsh", "bash", "fish"])
    
    with completion_tabs[0]:
        st.code("""
# 1. Download completion script
wget https://raw.githubusercontent.com/Equationzhao/g/master/completions/zsh/_g

# 2. Install to completion directory
mv _g ~/.zsh/completions  # or any directory in your $FPATH

# 3. Make sure these commands are in your .zshrc
autoload -Uz compinit
compinit
        """, language="bash")
    
    with completion_tabs[1]:
        st.code("""
# 1. Download completion script
wget https://raw.githubusercontent.com/Equationzhao/g/master/completions/bash/g-completion.bash

# 2. Add to .bashrc
echo "source /path/to/g-completion.bash" >> ~/.bashrc

# 3. Reload configuration
source ~/.bashrc
        """, language="bash")
    
    with completion_tabs[2]:
        st.code("""
# 1. Download completion script
wget https://raw.githubusercontent.com/Equationzhao/g/master/completions/fish/g.fish

# 2. Install to completion directory
mv g.fish ~/.config/fish/completions

# 3. Reload configuration
source ~/.config/fish/config.fish
        """, language="fish")

with shell_tabs[1]:
    alias_tabs = st.tabs(["bash", "zsh", "fish", "powershell", "nushell"])
    
    with alias_tabs[0]:
        st.code("""
# Add to .bashrc
eval "$(g --init bash)"

# Reload configuration
source ~/.bashrc
        """, language="bash")
    
    with alias_tabs[1]:
        st.code("""
# Add to .zshrc
eval "$(g --init zsh)"

# Reload configuration
source ~/.zshrc
        """, language="zsh")
    
    with alias_tabs[2]:
        st.code("""
# Add to fish configuration
g --init fish | source

# Reload configuration
source ~/.config/fish/config.fish
        """, language="fish")
    
    with alias_tabs[3]:
        st.code("""
# Add to your profile
Invoke-Expression (& { (g --init powershell | Out-String) })

# Find your profile path
echo $profile
        """, language="powershell")
    
    with alias_tabs[4]:
        st.code("""
# Add to $nu.env-path
^g --init nushell | save -f ~/.g.nu

# Add to $nu.config-path
source ~/.g.nu
        """, language="nu")

# Project Comparison
st.markdown('<div class="section-header" id="project-comparison">Project Comparison</div>', unsafe_allow_html=True)

st.markdown("""
g is highly inspired by the following excellent projects:
- [exa](https://github.com/ogham/exa) / [eza](https://github.com/eza-community/eza)
- [lsd](https://github.com/lsd-rs/lsd)
- [ls-go](https://github.com/acarl005/ls-go)
""")

# Create comparison table
comparison_data = {
    "Feature": ["Display Modes", "Unique Features", "Performance"],
    "eza": ["oneline, grid, across, tree, recurse", 
            "-Z: list security context, -@: list extended attributes and sizes", 
            "better"],
    "g": ["oneline, grid, across, zero, comma, table, json, markdown, tree, recurse", 
          "--mime: list mime types, --charset: list character sets", 
          "slower"]
}

df = pd.DataFrame(comparison_data)
st.table(df)


# Footer
st.markdown("---")
st.markdown("""
<div style="text-align: center">
    <p>g is an open-source project with MIT license</p>
    <p>© 2023-2025 Equationzhao</p>
</div>
""", unsafe_allow_html=True)