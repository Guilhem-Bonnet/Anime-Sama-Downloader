from utils.var import print_status, Colors
try:
    from moviepy import VideoFileClip
except ImportError:
    VideoFileClip = None
import subprocess
import sys
import shutil
import os

def print_status(message, status_type="info"):
    prefix = {
        "info": "[*]",
        "success": "[+]",
        "error": "[-]",
        "loading": "[...]"
    }.get(status_type.lower(), "[*]")
    print(f"{prefix} {message}")

def install_package(package_name):
    try:
        subprocess.check_call([sys.executable, "-m", "pip", "install", package_name])
        return True
    except subprocess.CalledProcessError:
        return False

def check_ffmpeg_installed():
    return shutil.which("ffmpeg") is not None

def convert_with_ffmpeg(ts_file, mp4_file):
    print_status(f"Converting {ts_file} to {mp4_file} using ffmpeg...", "loading")
    try:
        subprocess.check_call([
            "ffmpeg", "-y", "-i", ts_file,
            "-c:v", "libx264", "-preset", "ultrafast",
            "-c:a", "aac", mp4_file
        ])
        print_status("Conversion successful using ffmpeg!", "success")
        return True
    except subprocess.CalledProcessError as e:
        print_status(f"ffmpeg failed: {e}", "error")
        return False

def convert_with_moviepy(ts_file, mp4_file):
    try:
        from moviepy import VideoFileClip
    except ImportError:
        print_status("moviepy not installed.", "error")
        consent = input("Do you want to install moviepy? (y/n): ").lower()
        if consent == 'y':
            if not install_package("moviepy"):
                print_status("Failed to install moviepy.", "error")
                return False
        else:
            return False

    try:
        print_status(f"Converting {ts_file} to {mp4_file} using moviepy (slower)...", "loading")
        from moviepy import VideoFileClip
        clip = VideoFileClip(ts_file)
        clip.write_videofile(mp4_file, codec="libx264", audio_codec="aac", preset="superfast")
        clip.close()
        print_status("Conversion successful using moviepy!", "success")
        return True
    except Exception as e:
        print_status(f"MoviePy conversion failed: {e}", "error")
        return False
def convert_ts_to_mp4(input_path, output_path, pre_selected_tool=None):
    if not os.path.exists(input_path):
        print_status(f"Input file {input_path} does not exist", "error")
        return False, input_path
    if os.path.exists(output_path):
        print_status(f"Output file {output_path} already exists. deleting...", "error")
        try:
            os.remove(output_path)
        except Exception as e:
            print_status(f"Failed to delete existing output file: {e}", "error")
            return False, input_path
    if pre_selected_tool == 'ffmpeg':
    
        try:
            process = subprocess.Popen(
                ['ffmpeg', '-y', '-i', input_path, '-c:v', 'copy', '-c:a', 'copy', output_path],
                stdout=subprocess.PIPE,
                stderr=subprocess.STDOUT,
                text=True,
                bufsize=1,
                universal_newlines=True
            )

            for line in process.stdout:
                print(line, end='') 

            process.wait()

            if process.returncode == 0:
                print_status(f"Video converted successfully to {output_path}", "success")
                return True, output_path
            else:
                print_status("FFmpeg conversion failed", "error")
                return False, input_path
        except Exception as e:
            print_status(f"FFmpeg conversion failed: {str(e)}", "error")
            return False, input_path
    elif pre_selected_tool == 'moviepy':
        if not VideoFileClip:
            print_status("moviepy not installed, cannot convert to .mp4", "error")
            return False, input_path
        try:
            clip = VideoFileClip(input_path)
            clip.write_videofile(output_path, codec='libx264', audio_codec='aac')
            clip.close()
            return True, output_path
        except Exception as e:
            print_status(f"moviepy conversion failed: {str(e)}", "error")
            return False, input_path
    else:
        print_status("No valid conversion tool specified", "error")
        return False, input_path
 