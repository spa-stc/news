use std::{env, sync::LazyLock};

pub fn var(var: &str) -> Option<String> {
    const PREFIX: &str = "NEWSLETTER";

    env::var(format!("{}_{}", PREFIX, var)).ok()
}

static DEVELOPMENT: LazyLock<bool> = LazyLock::new(|| match var("DEVELOPMENT") {
    Some(val) => match val.to_lowercase().as_str() {
        "true" | "1" | "yes" | "y" => true,
        "false" | "0" | "no" | "n" => false,
        _ => false,
    },
    None => false,
});

pub fn is_development() -> bool {
    return *DEVELOPMENT;
}
