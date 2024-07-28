fn main() {
    let output = std::process::Command::new("git")
        .args(["rev-parse", "HEAD"])
        .output()
        .unwrap();

    let out = std::env::var("out").unwrap_or("/fake".into());

    let git_hash = String::from_utf8(output.stdout).unwrap();
    println!(
        "cargo:rustc-env=GIT_SHA={}",
        if git_hash.as_str() == "" {
            out
        } else {
            git_hash
        }
    );
}
