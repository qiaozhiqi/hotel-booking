-- 创建数据库
CREATE DATABASE IF NOT EXISTS hotel_booking DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE hotel_booking;

-- 用户表
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(100),
    phone VARCHAR(20),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 酒店表
CREATE TABLE IF NOT EXISTS hotels (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    address VARCHAR(255) NOT NULL,
    city VARCHAR(50) NOT NULL,
    description TEXT,
    rating DECIMAL(2,1) DEFAULT 0.0,
    image_url VARCHAR(255),
    price_range VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 房型表
CREATE TABLE IF NOT EXISTS rooms (
    id INT AUTO_INCREMENT PRIMARY KEY,
    hotel_id INT NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    price DECIMAL(10,2) NOT NULL,
    capacity INT NOT NULL,
    area INT,
    bed_type VARCHAR(50),
    amenities TEXT,
    image_url VARCHAR(255),
    total_count INT NOT NULL DEFAULT 1,
    available_count INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (hotel_id) REFERENCES hotels(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 订单表
CREATE TABLE IF NOT EXISTS orders (
    id INT AUTO_INCREMENT PRIMARY KEY,
    order_no VARCHAR(50) NOT NULL UNIQUE,
    user_id INT NOT NULL,
    hotel_id INT NOT NULL,
    room_id INT NOT NULL,
    check_in DATE NOT NULL,
    check_out DATE NOT NULL,
    guest_name VARCHAR(50) NOT NULL,
    guest_phone VARCHAR(20) NOT NULL,
    total_amount DECIMAL(10,2) NOT NULL,
    status ENUM('pending', 'confirmed', 'checked_in', 'checked_out', 'cancelled') DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (hotel_id) REFERENCES hotels(id),
    FOREIGN KEY (room_id) REFERENCES rooms(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 插入测试数据
-- 用户
INSERT INTO users (username, password, email, phone) VALUES 
('admin', 'admin123', 'admin@hotel.com', '13800138000'),
('testuser', 'test123', 'test@hotel.com', '13900139000');

-- 酒店
INSERT INTO hotels (name, address, city, description, rating, image_url, price_range) VALUES 
('上海外滩希尔顿酒店', '上海市黄浦区南京东路123号', '上海', '位于上海外滩的豪华五星级酒店，享有黄浦江美景。', 4.8, 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=luxury%20hotel%20shanghai%20bund%20modern%20architecture&image_size=landscape_4_3', '¥800-¥2000'),
('北京王府井洲际酒店', '北京市东城区王府井大街456号', '北京', '位于北京心脏地带的顶级商务酒店，交通便利。', 4.6, 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=luxury%20hotel%20beijing%20wangfujing%20business%20style&image_size=landscape_4_3', '¥600-¥1800'),
('广州天河希尔顿酒店', '广州市天河区天河路789号', '广州', '位于广州天河商务区的现代化酒店，设施完善。', 4.5, 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=modern%20hotel%20guangzhou%20tianhe%20business%20district&image_size=landscape_4_3', '¥500-¥1500'),
('深圳南山万豪酒店', '深圳市南山区科技园路101号', '深圳', '位于深圳科技园的商务酒店，紧邻各大科技企业。', 4.7, 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=business%20hotel%20shenzhen%20nanshan%20tech%20park&image_size=landscape_4_3', '¥700-¥1900'),
('杭州西湖四季酒店', '杭州市西湖区龙井路202号', '杭州', '坐落于杭州西湖畔的精品度假酒店，环境优美。', 4.9, 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=resort%20hotel%20hangzhou%20west%20lake%20scenic%20view&image_size=landscape_4_3', '¥1000-¥3000');

-- 房型
-- 上海酒店房型
INSERT INTO rooms (hotel_id, name, description, price, capacity, area, bed_type, amenities, image_url, total_count, available_count) VALUES 
(1, '豪华江景房', '宽敞明亮的豪华客房，享有黄浦江壮丽景色。配备舒适大床和现代化设施。', 1280.00, 2, 45, '大床', '免费WiFi, 空调, 电视, 迷你吧, 24小时热水', 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=luxury%20hotel%20room%20river%20view%20king%20bed%20modern&image_size=landscape_4_3', 20, 20),
(1, '行政套房', '高端行政套房，独立客厅和卧室，配备专属行政酒廊礼遇。', 2880.00, 2, 85, '大床', '免费WiFi, 空调, 电视, 迷你吧, 24小时热水, 行政酒廊', 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=luxury%20hotel%20suite%20living%20room%20bedroom%20elegant&image_size=landscape_4_3', 5, 5),
(1, '标准双床房', '经济实惠的标准客房，两张单人床，适合商务出行。', 880.00, 2, 35, '双床', '免费WiFi, 空调, 电视, 24小时热水', 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=standard%20hotel%20room%20twin%20beds%20clean%20simple&image_size=landscape_4_3', 30, 30);

-- 北京酒店房型
INSERT INTO rooms (hotel_id, name, description, price, capacity, area, bed_type, amenities, image_url, total_count, available_count) VALUES 
(2, '商务大床房', '专为商务旅客设计，配备宽大办公桌和高速网络。', 980.00, 2, 40, '大床', '免费WiFi, 空调, 电视, 办公桌, 24小时热水', 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=business%20hotel%20room%20work%20desk%20professional&image_size=landscape_4_3', 25, 25),
(2, '豪华套房', '宽敞的套房，独立起居空间，享受行政楼层特权。', 2580.00, 2, 80, '大床', '免费WiFi, 空调, 电视, 迷你吧, 24小时热水, 行政礼遇', 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=luxury%20hotel%20suite%20beijing%20style%20elegant&image_size=landscape_4_3', 8, 8);

-- 广州酒店房型
INSERT INTO rooms (hotel_id, name, description, price, capacity, area, bed_type, amenities, image_url, total_count, available_count) VALUES 
(3, '标准大床房', '舒适的标准客房，配备基本设施，性价比高。', 680.00, 2, 30, '大床', '免费WiFi, 空调, 电视, 24小时热水', 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=standard%20hotel%20room%20simple%20clean%20affordable&image_size=landscape_4_3', 40, 40),
(3, '行政大床房', '行政楼层客房，享有行政酒廊使用权，适合商务人士。', 1280.00, 2, 50, '大床', '免费WiFi, 空调, 电视, 迷你吧, 24小时热水, 行政酒廊', 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=executive%20hotel%20room%20guangzhou%20modern%20business&image_size=landscape_4_3', 15, 15);

-- 深圳酒店房型
INSERT INTO rooms (hotel_id, name, description, price, capacity, area, bed_type, amenities, image_url, total_count, available_count) VALUES 
(4, '科技大床房', '以科技为主题的客房，配备智能控制和高速网络。', 1080.00, 2, 42, '大床', '免费WiFi, 智能控制, 空调, 电视, 24小时热水', 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=tech%20themed%20hotel%20room%20smart%20devices%20modern&image_size=landscape_4_3', 20, 20),
(4, '总裁套房', '顶级套房，奢华装修，专属管家服务，彰显尊贵。', 3880.00, 2, 120, '大床', '免费WiFi, 智能控制, 空调, 电视, 迷你吧, 24小时管家服务', 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=presidential%20hotel%20suite%20luxury%20shenzhen%20elegant&image_size=landscape_4_3', 3, 3);

-- 杭州酒店房型
INSERT INTO rooms (hotel_id, name, description, price, capacity, area, bed_type, amenities, image_url, total_count, available_count) VALUES 
(5, '湖景豪华房', '享有西湖美景的豪华客房，装饰典雅，环境宁静。', 1580.00, 2, 55, '大床', '免费WiFi, 空调, 电视, 迷你吧, 24小时热水, 观景阳台', 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=lake%20view%20hotel%20room%20hangzhou%20west%20lake%20elegant&image_size=landscape_4_3', 18, 18),
(5, '别墅套房', '独立别墅套房，私人花园和泳池，尽享奢华度假体验。', 4880.00, 4, 200, '大床+双床', '免费WiFi, 空调, 电视, 迷你吧, 24小时管家, 私人泳池, 花园', 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=villa%20suite%20hotel%20room%20private%20pool%20garden%20luxury&image_size=landscape_4_3', 2, 2);
