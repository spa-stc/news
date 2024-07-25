use tracing::warn;

use crate::static_files::StaticFiles;

pub struct StaticFileFilter {
    files: StaticFiles,
}

impl StaticFileFilter {
    pub fn new(files: StaticFiles) -> Self {
        Self { files }
    }
}

impl tera::Filter for StaticFileFilter {
    fn filter(
        &self,
        value: &tera::Value,
        _args: &std::collections::HashMap<String, tera::Value>,
    ) -> tera::Result<tera::Value> {
        let key = value.as_str().unwrap_or_else(|| {
            warn!("Invalid input to static file filter");
            ""
        });
        let s = format!(
            "/static/{}",
            self.files.lookup_key(key).unwrap_or_else(|| {
                warn!("Static file with: {}, not found", key);
                "<invalid>"
            })
        );
        Ok(s.into())
    }
}
