# youtube_uploader

This is a lightweight python program that will take all videos in a directory and upload it to youtube.

## Setup
You will need to acquire a client secret from google following this doc:  
https://developers.google.com/identity/gsi/web/guides/get-google-api-clientid

In the api console, use the download as json button and save the secret somewhere
on your desktop. You will need its file path later.

## Running the code
This project builds from bazel, and I recommend using bazelisk (https://github.com/bazelbuild/bazelisk)
for ease of use.

To run this project's code, run:  
```
bazel run //youtube_uploader:main -- --client-secrets-path <downloaded json file path> --video-directory-root <root of video directory>
```