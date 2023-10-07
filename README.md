# pistasjis Attestator

**An open-source tool that allows you to check how apps on your system respect your privacy and security.**

## Install

You can grab a binary from [Releases](/releases). Make sure to verify the checksum before running.

## Build

Building the app is very, very simple. All you need on your system is [Go](https://go.dev) installed.

> [!NOTE]
> pistasjis Attestator does not currently work on Linux. Support will be added to the future.

Something like this should work:

```
git clone https://github.com/pistasjis/attestator
cd attestator
go build .
```

and you'll have a binary ready.

## Add app

Adding an app to the database is pretty easy. We use JSON and you can find the JSON at assets/apps.json.

For example, you can make an entry for KeePassXC with the reason "Open-source, local-only password manager with strict security" and verdict "good" like this:

```json
{
    "displayName": "KeePassXC",
    "reason": "Open-source, local-only password manager with strict security.",
    "verdict": "good"
}
```

Some guidelines for what's good, "meh" and bad will be posted at some point in the future.

## Privacy

The app only makes one request. The request is only for fetching the JSON (also referred to as "database") containing the list of apps from GitHub. If "raw.githubusercontent.com" is blocked on your network, the app will not work properly.