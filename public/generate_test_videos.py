# 导入必要的库
import os
import cv2
import numpy as np
import random
import multiprocessing
import subprocess
import time
from typing import List
from PIL import Image, ImageDraw, ImageFont

test_videos_directory = "test_videos"
audio_directory = "countdown_audio"
n_audio = len(os.listdir(audio_directory))


def generate_video(
    video_path,
    audio_path,
    text,
    duration,
    fps=30,
    size=(640, 480),
    bg_color=(54, 47, 0),  # b g r
    text_color=(227, 246, 253),
) -> bool:
    try:
        temp_video_path = video_path + ".temp.mp4"
        # 计算倒计时的初始值，假设每一秒减一
        countdown = duration
        # 计算视频的总帧数
        frames = (duration + 1) * fps

        # 创建一个图像对象，默认分辨率为 *size
        image = Image.fromarray(np.full((*size, 3), bg_color, dtype=np.uint8))

        # 创建一个绘制对象和字体对象
        draw = ImageDraw.Draw(image, mode="RGB")
        font = ImageFont.truetype("arial.ttf", 64)

        # 获取指定字符串和倒计时的宽度和高度
        text_width, text_height = draw.textsize(text, font=font)
        countdown_width, countdown_height = draw.textsize(str(countdown), font=font)

        # 计算指定字符串和倒计时的位置
        text_x = (image.width - text_width) // 2
        text_y = (image.height - text_height) // 4
        countdown_x = (image.width - countdown_width) // 2
        countdown_y = (image.height - countdown_height) * 3 // 4

        # 创建一个视频写入对象，格式为 MP4，默认帧率为 fps，分辨率为 *size
        video_writer = cv2.VideoWriter(
            temp_video_path,
            cv2.VideoWriter_fourcc(*"mp4v"),
            float(fps),
            (image.width, image.height),
            True,
        )

        # 使用一个循环来生成视频帧
        for i in range(frames):
            # 绘制指定字符串和倒计时在图像对象上，颜色为白色
            draw.text((text_x, text_y), text, fill=text_color, font=font)
            draw.text(
                (countdown_x, countdown_y), str(countdown), fill=text_color, font=font
            )

            # 将图像对象转换为 numpy 数组，并写入视频文件中
            frame = np.array(image)
            video_writer.write(frame)

            # 如果是每一秒的最后一帧，就更新倒计时的值
            if i % fps == fps - 1:
                # 如果倒计时已经结束，就跳出循环
                if countdown < 0:
                    break

                countdown -= 1

                # 否则，重新创建一个空白的图像对象和绘制对象
                image = Image.fromarray(np.full((*size, 3), bg_color, dtype=np.uint8))
                draw = ImageDraw.Draw(image, mode="RGB")

        # 关闭视频写入对象，并释放资源
        video_writer.release()

        # 使用FFmpeg合并音频和视频
        command = [
            "ffmpeg",
            "-y",
            "-i",
            temp_video_path,
            "-i",
            audio_path,
            "-c",
            "copy",
            "-strict",
            "experimental",
            "-map",
            "0:v",
            "-map",
            "1:a",
            video_path,
        ]

        # 执行FFmpeg命令
        subprocess.run(
            command, stdout=subprocess.PIPE, stderr=subprocess.PIPE, check=True
        )

        # 打印生成的信息
        print(f"完成生成 {video_path}，视频时长为 {duration} 秒")

        return True
    except subprocess.CalledProcessError as e:
        print("合并音频和视频失败")
        print(video_path, audio_path, text, duration)
        print("程序返回值：", e.returncode)
        print("程序输出：")
        print(e.output.decode("utf-8"))
        print("标准输出：")
        print(e.stdout.decode("utf-8"))
        print("标准错误：")
        print(e.stderr.decode("utf-8"))
        return False
    except Exception as e:
        print(e)
        return False
    finally:
        # 删除临时视频文件
        try:
            os.remove(temp_video_path)
        except FileNotFoundError:
            pass


# 定义一个函数，用于接收一个字符串和一个视频数量，然后生成对应的视频
def generate_videos(letters: str, num: int) -> int:
    video_dir = os.path.join(test_videos_directory, letters)
    try:
        os.mkdir(video_dir)
    except FileExistsError as e:
        print(e)

    num_successfully_generated = 0
    # 使用一个循环来生成多个视频
    for j in range(num):
        # 生成一个由字母和下标组成的 text
        text = f"{letters}{j + 1}"
        # 生成一个随机的时长，范围为 0 到 10
        duration = random.randint(0, n_audio - 1)
        video_path = os.path.join(video_dir, f"{j + 1}.mp4")
        audio_path = os.path.join(audio_directory, f"{duration}.mp3")

        # 调用函数来生成视频，并传入相应的参数
        if generate_video(video_path, audio_path, text, duration):
            num_successfully_generated += 1

    return num_successfully_generated


# 0 -> A 1 -> B 2 -> C ... 25 -> Z
# 26 -> AA 27 -> AB 28 -> AC ...
def num_to_alpha(num: int) -> str:
    if num < 26:
        return chr(num + ord("A"))
    else:
        return num_to_alpha(num // 26 - 1) + num_to_alpha(num % 26)


def main():
    # 创建一个空的进程池
    prc_pool = multiprocessing.Pool()
    # 创建一个空的列表，用于存储所有的结果
    results = []

    try:
        os.mkdir(test_videos_directory)
    except FileExistsError as e:
        print(e)

    # 使用一个循环来遍历所有字母
    for i in range(32):
        # 生成一个随机的视频数量，范围为 1 到 10
        num = random.randint(1, 10)

        letters = num_to_alpha(i)

        # 创建一个进程对象，并传入相应的参数
        result = prc_pool.apply_async(generate_videos, args=(letters, num))

        # 将进程对象添加到列表中
        results.append(result)

    # 关闭进程池
    prc_pool.close()

    # 等待进程结束
    sum_successfully_generated = 0
    for result in results:
        result.wait()
        if result.successful():
            sum_successfully_generated += result.get()
        print("进程", os.getpid(), "结束")

    # 等待所有进程结束
    prc_pool.join()

    print("成功生成", sum_successfully_generated, "个视频")


multiprocessing.freeze_support()

if __name__ == "__main__":
    main()
