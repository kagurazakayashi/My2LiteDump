import mysql.connector
from mysql.connector import errorcode
from faker import Faker
import random
import uuid

# === 設定資料庫連線參數 ===
DB_CONFIG = {
    "host": "127.0.0.1",
    "user": "test1",
    "password": "test1",
    "database": "test1"
}

def create_connection():
    """ 創建並返回 MySQL 資料庫連線。 """
    try:
        conn = mysql.connector.connect(**DB_CONFIG)
        print("成功連線到 MySQL 資料庫")
        return conn
    except mysql.connector.Error as err:
        print(f"資料庫連線失敗: {err}")
        exit(1)

def table_exists(conn, table_name="contacts"):
    """ 檢查資料表是否存在。 """
    cursor = conn.cursor()
    cursor.execute(f"SHOW TABLES LIKE '{table_name}'")
    result = cursor.fetchone()
    cursor.close()
    return result is not None

def create_table(conn):
    """ 創建 `contacts` 資料表，包含 LGBTQ+ 友好的性別選項與 UUID。 """
    cursor = conn.cursor()
    create_table_sql = """
    CREATE TABLE contacts (
        id INT AUTO_INCREMENT PRIMARY KEY COMMENT '自增主鍵',
        employee_id BIGINT UNIQUE NOT NULL COMMENT '工號（10 位數字）',
        uuid VARCHAR(36) UNIQUE NOT NULL COMMENT '唯一 UUID',
        name VARCHAR(50) NOT NULL COMMENT '姓名',
        english_name VARCHAR(50) COMMENT '英文名',
        gender ENUM(
            '男性', '女性', '非二元', '跨性別男性', '跨性別女性', 
            '性別流動', '雙性', '酷兒', '其他'
        ) NOT NULL COMMENT '性別',
        department VARCHAR(50) COMMENT '部門',
        phone VARCHAR(20) UNIQUE COMMENT '電話號碼（唯一）',
        birthdate DATE COMMENT '生日',
        hire_date DATETIME COMMENT '入職時間',
        employee_level TINYINT CHECK (employee_level BETWEEN 1 AND 10) COMMENT '員工等級（1-10）',
        id_card VARCHAR(18) UNIQUE NOT NULL COMMENT '身分證號（唯一）',
        is_active BOOLEAN DEFAULT TRUE COMMENT '是否在職（1: 在職，0: 離職）'
    ) COMMENT='員工通訊錄表';
    """
    print(f"執行 SQL:\n{create_table_sql}\n")
    try:
        cursor.execute(create_table_sql)
        conn.commit()
        print("資料表 contacts 創建成功（含欄位註釋）")
    except mysql.connector.Error as err:
        print(f"創建資料表失敗: {err}")
    finally:
        cursor.close()

def generate_sample_data(n=10000):
    """ 生成 n 筆隨機員工資料，確保工號、UUID 和身分證號唯一。 """
    fake = Faker("zh_TW")
    fake_en = Faker("en_US")
    data = []
    genders = ["男性", "女性"]
    tgenders = ["非二元", "跨性別男性", "跨性別女性", "性別流動", "雙性", "酷兒", "其他"]
    departments = ["研發部", "市場部", "銷售部", "人事部", "財務部", "運營部", "行政部"]

    existing_employee_ids = set()
    existing_id_cards = set()

    while len(data) < n:
        employee_id = random.randint(1000000000, 9999999999)
        id_card = fake.ssn()
        unique_uuid = str(uuid.uuid4())

        if employee_id in existing_employee_ids or id_card in existing_id_cards:
            continue  # 確保 employee_id 和 id_card 唯一

        # 80% 機率為「男性」或「女性」
        gender = random.choices(
            population=genders + tgenders,
            weights=[40, 40] + [20 / len(tgenders)] * len(tgenders),
            k=1
        )[0]

        name = fake.name()
        english_name = fake_en.first_name() + " " + fake_en.last_name()
        department = random.choice(departments)
        phone = fake.phone_number()
        birthdate = fake.date_of_birth(minimum_age=20, maximum_age=60)
        hire_date = fake.date_time_between(start_date="-25y", end_date="now")
        employee_level = random.randint(1, 10)
        is_active = random.choice([True, False])

        data.append((employee_id, unique_uuid, name, english_name, gender, department, phone, birthdate, hire_date, employee_level, id_card, is_active))

    return data

def insert_data(conn, data):
    """ 每 10 條數據執行一次批量插入，並輸出所有 SQL。 """
    cursor = conn.cursor()
    insert_sql = """
    INSERT INTO contacts (employee_id, uuid, name, english_name, gender, department, phone, birthdate, hire_date, employee_level, id_card, is_active)
    VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s)
    ON DUPLICATE KEY UPDATE
        name = VALUES(name),
        english_name = VALUES(english_name),
        gender = VALUES(gender),
        department = VALUES(department),
        phone = VALUES(phone),
        birthdate = VALUES(birthdate),
        hire_date = VALUES(hire_date),
        employee_level = VALUES(employee_level),
        id_card = VALUES(id_card),
        is_active = VALUES(is_active);
    """

    batch_size = 10
    for i in range(0, len(data), batch_size):
        batch = data[i:i + batch_size]
        print(f"執行 SQL（第 {i // batch_size + 1} 批）:\n{insert_sql}\n")
        print(f"示例數據: {batch[0]}\n")
        try:
            cursor.executemany(insert_sql, batch)
            conn.commit()
        except mysql.connector.Error as err:
            print(f"插入數據失敗: {err}")

    cursor.close()

if __name__ == "__main__":
    conn = create_connection()

    if not table_exists(conn, "contacts"):
        create_table(conn)
        print("創建新資料表並插入資料")
    else:
        print("資料表已存在，直接插入資料")

    sample_data = generate_sample_data(10000)
    insert_data(conn, sample_data)

    conn.close()
    print("資料庫操作完成")
