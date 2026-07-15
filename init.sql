-- ========================================================
-- 1. CLEANUP (Optional: Drops tables if they already exist)
-- ========================================================
DROP TABLE IF EXISTS applications;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS application_gates;
DROP TABLE IF EXISTS user_types;

-- ========================================================
-- 2. SCHEMA DEFINITION (With Auto-Incrementing IDs)
-- ========================================================

-- Table 1: User Types (Reference Table)
CREATE TABLE user_types (
    user_type_id INT AUTO_INCREMENT PRIMARY KEY,
    type_name VARCHAR(50) NOT NULL
);

-- Table 2: Users (Links to UserTypes)
CREATE TABLE users (
    user_id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    user_type_id INT,
    CONSTRAINT fk_user_type FOREIGN KEY (user_type_id) REFERENCES user_types(user_type_id)
);

-- Table 3: Application Gates (Timeline/Enrollment Records)
CREATE TABLE application_gates (
    gate_id INT AUTO_INCREMENT PRIMARY KEY,
    gate_name VARCHAR(100) NOT NULL,
    active_year INT NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL
);

-- Table 4: Applications (Links to Users and Gates)
CREATE TABLE applications (
    application_id INT AUTO_INCREMENT PRIMARY KEY,
    application_title VARCHAR(200) NOT NULL,
    application_status VARCHAR(50) DEFAULT 'Submitted',
    submission_date DATE NOT NULL,
    user_id INT,
    gate_id INT,
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(user_id),
    CONSTRAINT fk_gate FOREIGN KEY (gate_id) REFERENCES application_gates(gate_id)
);

-- ========================================================
-- 3. DATA INSERTION (20 Rows Per Table)
-- ========================================================

-- Populate user_types
INSERT INTO user_types (type_name) VALUES 
('Undergraduate'), ('Graduate'), ('Doctoral'), ('Faculty'), ('Staff'), 
('Alumni'), ('Admin'), ('Guest'), ('Transfer'), ('International'), 
('Exchange'), ('Post-Doc'), ('Researcher'), ('Non-Degree'), ('Visiting'), 
('Continuing Ed'), ('High School Dual'), ('Veteran'), ('Scholarship-Only'), ('Parent');

-- Populate users
INSERT INTO users (username, email, user_type_id) VALUES 
('j_smith', 'j.smith@univ.edu', 1), ('m_garcia', 'm.garcia@univ.edu', 1), 
('a_chen', 'a.chen@univ.edu', 2), ('b_wilson', 'b.wilson@univ.edu', 7), 
('k_patel', 'k.patel@univ.edu', 1), ('l_thompson', 'l.thompson@univ.edu', 1), 
('s_jones', 's.jones@univ.edu', 9), ('r_miller', 'r.miller@univ.edu', 3), 
('e_davis', 'e.davis@univ.edu', 10), ('h_white', 'h.white@univ.edu', 1), 
('t_brown', 't.brown@univ.edu', 1), ('p_clark', 'p.clark@univ.edu', 1), 
('c_lewis', 'c.lewis@univ.edu', 12), ('m_walker', 'm.walker@univ.edu', 1), 
('d_hall', 'd.hall@univ.edu', 1), ('j_young', 'j.young@univ.edu', 1), 
('n_king', 'n.king@univ.edu', 1), ('k_wright', 'k.wright@univ.edu', 1), 
('g_scott', 'g.scott@univ.edu', 1), ('f_green', 'f.green@univ.edu', 1);

-- Populate application_gates (Enrollment Periods)
INSERT INTO application_gates (gate_name, active_year, start_date, end_date) VALUES 
('Spring 2026 Early', 2026, '2025-09-01', '2025-10-15'),
('Spring 2026 Regular', 2026, '2025-10-16', '2025-12-01'),
('Summer 2026 Session I', 2026, '2026-01-01', '2026-03-01'),
('Summer 2026 Session II', 2026, '2026-02-01', '2026-04-01'),
('Fall 2026 Early Action', 2026, '2025-11-01', '2026-01-15'),
('Fall 2026 Regular Decision', 2026, '2026-01-16', '2026-05-01'),
('Fall 2026 Late Enrollment', 2026, '2026-05-02', '2026-08-15'),
('Winter 2026 Short Term', 2026, '2026-09-01', '2026-11-01'),
('Spring 2027 Early', 2027, '2026-09-01', '2026-10-15'),
('Spring 2027 Regular', 2027, '2026-10-16', '2026-12-01'),
('Scholarship 2026-A', 2026, '2025-12-01', '2026-02-01'),
('Scholarship 2026-B', 2026, '2026-03-01', '2026-05-01'),
('Transfer Fall 2026', 2026, '2026-02-01', '2026-06-01'),
('International Fall 2026', 2026, '2025-10-01', '2026-02-01'),
('Housing Window 2026', 2026, '2026-03-01', '2026-07-01'),
('Masters Research Fall', 2026, '2025-09-15', '2025-12-15'),
('Study Abroad Spring', 2026, '2025-05-01', '2025-09-01'),
('Study Abroad Fall', 2026, '2025-12-01', '2026-03-01'),
('Continuing Ed Q1', 2026, '2025-12-01', '2026-01-10'),
('Continuing Ed Q2', 2026, '2026-03-01', '2026-04-10');

-- Populate applications (Mapping users to gates)
INSERT INTO applications (application_title, application_status, user_id, gate_id, submission_date) VALUES 
('Computer Science Major', 'Accepted', 1, 1, '2025-09-15'),
('Mechanical Engineering', 'Pending', 2, 1, '2025-09-20'),
('MBA Program', 'Under Review', 3, 16, '2025-10-20'),
('Admin System Access', 'Approved', 4, 15, '2026-01-12'),
('Physics Research', 'Submitted', 5, 3, '2026-01-05'),
('Liberal Arts Regular', 'Pending', 6, 6, '2026-01-20'),
('Transfer Art History', 'Under Review', 7, 13, '2026-02-15'),
('PhD Dissertation Fund', 'Pending', 8, 11, '2026-01-10'),
('Global Studies Intl', 'Accepted', 9, 14, '2025-10-10'),
('Biology Early Dec', 'Withdrawn', 10, 5, '2026-01-10'),
('Late Enrollment Chem', 'Under Review', 11, 7, '2026-05-10'),
('General Studies', 'Accepted', 12, 2, '2025-11-01'),
('Post-Doc Fellowship', 'Submitted', 13, 16, '2025-11-15'),
('Freshman Housing', 'Approved', 14, 15, '2026-03-10'),
('Mathematics Major', 'Pending', 15, 6, '2026-03-15'),
('History Major', 'Accepted', 16, 6, '2026-04-01'),
('Sociology Transfer', 'Pending', 17, 13, '2026-04-20'),
('Psychology Major', 'Submitted', 18, 6, '2026-02-10'),
('Music Performance', 'Accepted', 19, 1, '2025-09-30'),
('Chemistry 101 Cont Ed', 'Accepted', 20, 19, '2026-01-02');
