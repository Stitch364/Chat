# 创建表

# 用户表
create table if not exists users(
    id bigint AUTO_INCREMENT not null primary key,     # ID
    email varchar(255) not null unique, # 邮箱
    password varchar(255) not null,     # 密码
    created_at timestamp not null default now()    # 创建时间
);

# 账号
create table if not exists accounts(
    id bigint not null primary key,     # ID
    user_id bigint not null references users (id) on delete cascade on update cascade,# 用户 ID(外键)
    name varchar(255) not null,         # 账号名
    avatar varchar(255) not null,       # 头像
    gender enum('男','女','未知') not null default '未知',    # 性别
    signature varchar(1000) not null default '这个用户很懒，什么也没有留下~',  # 个性签名
    created_at timestamp not null default now(),    # 创建时间
    constraint unique (user_id, name) #同一个用户可以拥有多个账号，但是账号名不能重复
);
-- 账号名和头像索引(可以加快根据账号名和头像查询的速度)
create index account_index_name_avatar on accounts(name, avatar);

# 群组或好友
create table if not exists relations(
    id bigint AUTO_INCREMENT not null primary key, -- id
    relation enum('group','friend') not null, -- 关系类型

    name varchar(50), -- 群名称
    description varchar(255), -- 群描述
    avatar varchar(255), -- 群头像

    account1_id bigint, -- 好友 1 的账号 id
    account2_id bigint, -- 好友 2 的账号 id

    created_at timestamp not null default now()   # 创建时间

    # 群信息 和 好友信息只能存在一种
);

# 账号对群组或好友关系的设置
create table settings(
    account_id bigint not null references accounts (id) on delete cascade on update cascade, -- 账号id（外键）
    relation_id bigint not null references relations (id) on delete cascade on update cascade, -- 关系 id（外键）
    nick_name varchar(255) not null,                -- 备注（默认是账户名或群组名）
    is_not_disturb boolean not null default false,  -- 是否免打扰
    is_pin boolean not null default false,      -- 是否置顶
    pin_time timestamp not null default now(),  -- 置顶时间
    is_show boolean not null default true,      -- 是否显示
    last_show timestamp not null default now(), -- 最后一次显示时间
    is_leader boolean not null default false, -- 是否是群主，仅对群组有效
    is_self boolean not null default false -- 是否是自己对自己的关系，仅对好友有效

);

-- 昵称索引(可以加快根据昵称查询的速度)
create index relation_setting_nickname on settings (nick_name);
-- 账户ID和关系ID的复合索引
create index setting_idx_account_id_relation_id on settings (account_id, relation_id);

-- 好友申请
create table if not exists applications(
     account1_id bigint not null references accounts (id) on delete cascade on update cascade, -- 申请者账号 id（外键）
     account2_id bigint not null references accounts (id) on delete cascade on update cascade, -- 被申请者账号 id（外键）
     apply_msg text not null, -- 申请信息
     refuse_msg text not null, -- 拒绝信息
     status enum('已申请','已同意','已拒绝','等待验证') not null default '已申请', -- 申请状态
     create_at timestamp not null default now(), -- 创建时间
     update_at timestamp not null default now(), -- 更新时间
     constraint f_a_pk primary key (account1_id, account2_id,create_at)
);

-- 文件记录
create table if not exists files
(
    id bigint AUTO_INCREMENT not null primary key, -- 文件 id
    file_name varchar(255) not null, -- 文件名称
    file_type enum('image','file') not null, -- 文件类型
    file_size bigint not null, -- 文件大小 byte
    file_key varchar(255) not null, -- 文件 key 用于 oss 中删除文件
    url varchar(255) not null, -- 文件 url
    relation_id bigint references relations (id) on delete cascade on update cascade, -- 关系 id（外键）（群组/好友）
    account_id bigint references accounts (id) on delete cascade on update cascade, -- 发送账号 id（外键）
    create_at timestamp not null default now() -- 创建时间
);

-- 文件关系id索引
create index file_relation_id on files (relation_id);

-- 消息
create table if not exists messages
(
    id bigint AUTO_INCREMENT primary key, -- 消息 id
    notify_type enum('system', 'common') not null, -- 消息通知类型 system:系统消息，common:普通消息
    msg_type enum('text', 'file')not null, -- 消息类型 text:文本消息，file:文件消息
    msg_content text not null, -- 消息内容
    msg_extend json, -- 消息扩展信息
    file_id bigint references files (id) on delete cascade on update cascade, -- 文件 id（外键），如果不是文件类型则为 null
    account_id bigint references accounts (id) on delete set null on update cascade, -- 发送账号 id（外键）
    rly_msg_id bigint references messages (id) on delete cascade on update cascade, -- 回复消息 id，没有则为 null（外键）
    relation_id bigint not null references relations (id) on delete cascade on update cascade, -- 关系 id（外键）
    create_at timestamp not null default now(), -- 创建时间
    is_revoke boolean not null default false, -- 是否撤回
    is_top boolean not null default false, -- 是否置顶
    is_pin boolean not null default false, -- 是否pin(将消息固定在会话顶部)
    pin_time timestamp not null default now(), -- pin时间
    read_ids json, -- 已读用户 id 集合 默认是空的json数组
    #msg_content_tsy tsvector, -- 消息分词
    is_delete int not null default 0,
    check (notify_type = 'common' or notify_type = 'system'), -- 系统消息时发送账号 id 为 null
    check (msg_type = 'text' or msg_type = 'file') -- 文件消息时文件 id 不能为 null
);
-- 创建时间索引
create index msg_create_at on messages (create_at);

# -- 创建将已读设置为null的触发器
# CREATE TRIGGER  before_insert_messages
#     BEFORE INSERT ON messages
#     FOR EACH ROW
# BEGIN
#     IF NEW.read_ids IS NULL THEN
#         SET NEW.read_ids = JSON_ARRAY();
#     END IF;
# END;

-- 群通知
create table if not exists group_notify
(
    id bigint AUTO_INCREMENT primary key, -- 群通知 id
    relation_id bigint references relations (id) on delete cascade on update cascade, -- 关系 id（外键）
    msg_content text not null, -- 消息内容
    msg_expand json, -- 消息扩展信息
    account_id bigint references accounts (id) on delete cascade on update cascade, -- 发送账号 id（外键）
    create_at timestamp not null default now(), -- 创建时间
    read_ids json -- 已读用户 id 集合
    #msg_content_tsv tsvector -- 消息分词
);

# -- 创建将群通知已读设置为null的触发器
# CREATE TRIGGER before_insert_group_notify
#     BEFORE INSERT ON group_notify
#     FOR EACH ROW
# BEGIN
#     IF NEW.read_ids IS NULL THEN
#         SET NEW.read_ids = JSON_ARRAY();
#     END IF;
# END;
