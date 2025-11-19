use crate::database::DataBase;

mod connection;
mod database;

const MAX_CONN: i8 = 8;

pub fn init_database(dsn: String) -> DataBase {
    let mut db = DataBase::new(&dsn);
    db.connect(&dsn, MAX_CONN);
    db
}
