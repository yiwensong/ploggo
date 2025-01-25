import argparse
import asyncio
import datetime
import glob
import os

import yaml
from apiclient.errors import HttpError

import youtube_uploader.uploader as uploader


def get_arg_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser(
        description="youtube uploader for outplayed",
    )
    parser.add_argument(
        "--client-secrets-path",
        required=True,
        help="path to client secrets file",
    )
    parser.add_argument(
        "--video-directory-root",
        required=True,
        help="root of the directory where videos are stored",
    )
    parser.add_argument(
        "--run-env",
        default=os.path.join(os.environ.get("HOME"), "youtube-uploader"),
        help="a directory in which this app will store state",
    )
    parser.add_argument(
        "--upload-after",
        default=None,
        help="if set, only videos on or after this date will be included in upload",
    )
    return parser


PREV_UPLOADED_YAML = "previously_uploaded.yaml"
VIDEO_CATEGORY = 22


def get_videos_to_upload(
    run_env: str, video_directory_root: str, upload_after: datetime.datetime
) -> list[str]:
    """returns all videos that need to be uploaded"""
    prev_uploaded_files = set()
    previously_uploaded_yaml = os.path.join(run_env, PREV_UPLOADED_YAML)
    
    # in case this is the first time we run this tool, we want
    # to do the setup and not fail.
    os.makedirs(run_env)
    if not os.path.exists(previously_uploaded_yaml):
        with open(previously_uploaded_yaml, "w") as new_file:
            new_file.write("")

    with open(previously_uploaded_yaml, "r") as prev_file:
        prev_uploaded_files = set(yaml.safe_load(prev_file))

    videos = glob.glob(os.path.join(video_directory_root, "**.mp4"))
    upload_after_ts = upload_after.timestamp()

    for video in videos:
        if os.stat(video).st_mtime < upload_after_ts:
            continue
        if video in prev_uploaded_files:
            continue


def write_uploaded(run_env: str, newly_uploaded: list[str]) -> None:
    previously_uploaded_yaml = os.path.join(run_env, PREV_UPLOADED_YAML)
    with open(previously_uploaded_yaml, "r+") as prev_file:
        prev_uploaded_files = yaml.safe_load(prev_file)
        prev_uploaded_files += newly_uploaded
        prev_file.seek(0)
        yaml.dump(prev_uploaded_files, prev_file)


async def main_async() -> None:
    args = get_arg_parser().parse()

    upload_after = (
        datetime.datetime.today() - datetime.timedelta(days=1)
        if args.upload_after is None
        else datetime.datetime.fromisoformat(args.upload_after)
    )
    videos = get_videos_to_upload(
        args.run_env,
        args.video_directory_root,
        upload_after,
    )

    youtube = uploader.get_authenticated_service(args.client_secrets_path)

    async with asyncio.TaskGroup() as tg:
        for video in videos:
            upload_options = uploader.UploadOptions(
                file=video,
                title=os.path.basename(video),
                description="description",
                category=VIDEO_CATEGORY,
                privacy_status="unlisted",
                keywords="",
            )
            tg.create_task(uploader.initialize_upload_async(youtube, upload_options))
            
        await tg
    
    write_uploaded(args.run_env, videos)


if __name__ == "__main__":
    asyncio.run(main_async())
