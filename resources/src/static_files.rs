use std::{
    collections::HashMap,
    fmt::Debug,
    fs::File,
    io::{BufReader, Read},
    path::Path,
};

use bytes::Bytes;
use sha2::{Digest, Sha256};
use tracing::info;

#[derive(Clone, Default)]
pub struct StaticFiles {
    by_key: HashMap<String, String>,
    files: HashMap<String, (Bytes, &'static str)>,
}

impl StaticFiles {
    pub fn register_dir<P: AsRef<Path> + Debug>(&mut self, path: P) -> Result<(), crate::Error> {
        let path = path.as_ref();

        info!("building statics in: {:?}", path);
        for file in std::fs::read_dir(path)? {
            let file = file?;

            if file.file_type()?.is_file() {
                let file = file.file_name();
                let name = Path::new(&file);
                let ext = name.extension().unwrap_or_default().to_string_lossy();
                self.register_file(&file.to_string_lossy(), &ext, &path.join(name))?;
            }
        }

        Ok(())
    }

    pub fn register(&mut self, key: &str, extension: &str, buf: &[u8]) {
        let mime =
            mime_type_from(extension, buf).expect(&format!("mime type was not known for {}", key));

        let mut hash = Sha256::default();
        hash.update(buf);
        let hash: &[u8] = &hash.finalize();

        self.files.insert(
            to_hash_key(hash) + "." + extension,
            (Bytes::from(buf.to_vec()), mime),
        );
        self.by_key
            .insert(key.into(), to_hash_key(hash) + "." + extension);

        tracing::info!(
            "registered: '{}.{}' with hash: {}",
            key,
            extension,
            to_hash_key(hash)
        );
    }

    pub fn register_file<P: AsRef<Path>>(
        &mut self,
        key: &str,
        extension: &str,
        path: P,
    ) -> Result<(), crate::Error> {
        let mut reader = BufReader::new(File::open(&path)?);
        let mut buf = Vec::with_capacity(1024);
        reader.read_to_end(&mut buf)?;
        self.register(key, extension, &buf);

        Ok(())
    }

    pub fn lookup_key(&self, key: &str) -> Option<&str> {
        self.by_key.get(key).map(|x| x.as_str())
    }

    pub fn get_bytes_from_key(&self, key: &str) -> Option<(Bytes, &'static str)> {
        self.files.get(key).map(|x| (x.0.clone(), x.1))
    }
}

fn mime_type_from(extension: &str, buf: &[u8]) -> Option<&'static str> {
    match extension {
        "txt" => Some("text/plain"),
        "css" => Some("text/css"),
        "js" => Some("application/javascript"),
        _ => infer::get(buf).map(|x| x.mime_type()),
    }
}

fn to_hash_key(bytes: &[u8]) -> String {
    let mut s = "v0-".to_owned();
    for byte in bytes {
        s += &format!("{:02x}", byte);
    }
    s
}
