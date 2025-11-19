// database.rs
use crate::connection::Conn;

pub struct DataBase {
    pub(crate) dsn: String,
    conns: Vec<Conn>,
}

impl DataBase {
    pub fn new(dsn: &str) -> Self {
        DataBase {
            dsn: dsn.to_string(),
            conns: Vec::new(),
        }
    }

    pub fn connect(&mut self, dsn: &str, max_conn: i8) {
        self.dsn = dsn.to_string();
        for _ in 0..max_conn {
            let new_conn = Conn::init(max_conn);
            self.conns.push(new_conn);
        }
    }
}