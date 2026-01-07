"""
Path input with tab completion and favorites support
"""

import os
try:
    import readline
    HAS_READLINE = True
except ImportError:
    HAS_READLINE = False
    
import glob
from utils.var import Colors, print_status
from utils.config import get_favorite_paths, add_favorite_path, get_last_used_path, set_last_used_path

class PathCompleter:
    """Tab completion for file paths"""
    
    def __init__(self):
        self.matches = []
    
    def complete(self, text, state):
        """Return the next possible completion for 'text'."""
        if state == 0:
            # Expand ~ and environment variables
            text = os.path.expanduser(text)
            text = os.path.expandvars(text)
            
            # If empty or ends with /, complete directory contents
            if not text or text.endswith('/'):
                self.matches = glob.glob(text + '*')
            else:
                # Complete partial paths
                self.matches = glob.glob(text + '*')
            
            # Add trailing slash to directories
            self.matches = [
                (m + '/' if os.path.isdir(m) else m) 
                for m in self.matches
            ]
        
        try:
            return self.matches[state]
        except IndexError:
            return None

def setup_readline():
    """Setup readline for tab completion"""
    if not HAS_READLINE:
        return False
        
    try:
        # Set up tab completion
        completer = PathCompleter()
        readline.set_completer(completer.complete)
        readline.parse_and_bind('tab: complete')
        
        # Set completer delimiters (don't break on /)
        readline.set_completer_delims(' \t\n;')
        return True
    except Exception as e:
        return False

def get_save_directory_interactive():
    """
    Interactive directory selection with favorites and tab completion.
    Returns absolute path to save directory.
    """
    from utils.var import print_separator
    
    print(f"\n{Colors.BOLD}{Colors.HEADER}📁 DOWNLOAD LOCATION{Colors.ENDC}")
    print_separator()
    
    # Show favorites if any
    favorites = get_favorite_paths()
    last_used = get_last_used_path()
    
    if favorites or last_used:
        print(f"{Colors.BOLD}Quick options:{Colors.ENDC}")
        
        if last_used and os.path.exists(last_used):
            print(f"  {Colors.OKGREEN}L{Colors.ENDC}. Last used: {last_used}")
        
        for i, fav in enumerate(favorites, 1):
            if os.path.exists(fav):
                marker = f"{Colors.OKCYAN}{i}{Colors.ENDC}"
                print(f"  {marker}. {fav}")
        
        print(f"  {Colors.BOLD}C{Colors.ENDC}. Custom path (with tab completion)")
        print()
    
    # Setup tab completion
    has_completion = setup_readline()
    if has_completion:
        print_status("💡 Tip: Use TAB for path autocompletion", "info")
    
    while True:
        prompt = f"{Colors.BOLD}Choose location (L/1-{len(favorites)}/C or path): {Colors.ENDC}"
        choice = input(prompt).strip()
        
        # Quick selections
        if choice.upper() == 'L' and last_used and os.path.exists(last_used):
            save_dir = last_used
            break
        elif choice.isdigit() and 1 <= int(choice) <= len(favorites):
            fav_index = int(choice) - 1
            if os.path.exists(favorites[fav_index]):
                save_dir = favorites[fav_index]
                break
            else:
                print_status(f"Favorite path no longer exists: {favorites[fav_index]}", "error")
                continue
        elif choice.upper() == 'C' or choice == '':
            # Custom path input with tab completion
            path_prompt = f"{Colors.BOLD}Enter path (use TAB to complete): {Colors.ENDC}"
            custom_path = input(path_prompt).strip()
            
            if not custom_path:
                print_status("Path cannot be empty", "error")
                continue
            
            # Process the custom path
            custom_path = os.path.expanduser(custom_path)
            custom_path = os.path.expandvars(custom_path)
            custom_path = os.path.abspath(custom_path)
            
            # Create if doesn't exist
            if not os.path.exists(custom_path):
                create = input(f"{Colors.BOLD}Directory doesn't exist. Create it? (Y/n): {Colors.ENDC}").strip().lower()
                if create != 'n':
                    try:
                        os.makedirs(custom_path, exist_ok=True)
                        print_status(f"Created directory: {custom_path}", "success")
                    except Exception as e:
                        print_status(f"Failed to create directory: {e}", "error")
                        continue
                else:
                    continue
            
            # Ask to save as favorite
            if custom_path not in favorites:
                save_fav = input(f"{Colors.BOLD}Save as favorite? (Y/n): {Colors.ENDC}").strip().lower()
                if save_fav != 'n':
                    add_favorite_path(custom_path)
                    print_status("Added to favorites", "success")
            
            save_dir = custom_path
            break
        else:
            # Treat as direct path input
            custom_path = choice
            
            # Expand ~ and variables
            custom_path = os.path.expanduser(custom_path)
            custom_path = os.path.expandvars(custom_path)
            custom_path = os.path.abspath(custom_path)
            
            # Create if doesn't exist
            if not os.path.exists(custom_path):
                create = input(f"{Colors.BOLD}Directory doesn't exist. Create it? (Y/n): {Colors.ENDC}").strip().lower()
                if create != 'n':
                    try:
                        os.makedirs(custom_path, exist_ok=True)
                        print_status(f"Created directory: {custom_path}", "success")
                    except Exception as e:
                        print_status(f"Failed to create directory: {e}", "error")
                        continue
                else:
                    continue
            
            # Ask to save as favorite
            if custom_path not in favorites:
                save_fav = input(f"{Colors.BOLD}Save as favorite? (Y/n): {Colors.ENDC}").strip().lower()
                if save_fav != 'n':
                    add_favorite_path(custom_path)
                    print_status("Added to favorites", "success")
            
            save_dir = custom_path
            break
    
    # Save as last used
    set_last_used_path(save_dir)
    print_status(f"Saving to: {save_dir}", "info")
    
    return save_dir
