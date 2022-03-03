recipe "brew" "local" {
    commands = [
        "brew update",
        "brew upgrade",
        "brew upgrade k6",
    ]
}