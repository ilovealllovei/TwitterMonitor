create table channels
(
    id              varchar(255)         not null
        primary key,
    ownerId         int                  not null,
    isVerified      tinyint(1) default 0 null,
    name            varchar(255)         not null,
    description     text                 null,
    avatar          varchar(255)         null,
    chatLink        varchar(255)         null,
    isPublic        tinyint(1) default 0 null,
    isHot           tinyint(1) default 0 null,
    hotExpireAt     varchar(255)         null,
    createdAt       bigint               null,
    updatedAt       bigint               null,
    watchlist       json                 null,
    eventlist       json                 null,
    followerCount   varchar(255)         null,
    recentFollowers json                 null
);

create table follows
(
    id        varchar(36) not null
        primary key,
    userId    int         not null,
    channelId varchar(36) not null,
    createdAt bigint      not null
);

create index idx_user_channel
    on follows (userId, channelId);

create table twitter_info
(
    id         int auto_increment comment '自增主键，用于唯一标识每条记录'
        primary key,
    tweetsId   varchar(255) not null,
    twitterId  varchar(255) not null comment '推特id',
    content    longtext     null comment '内容',
    chainId    varchar(255) null comment '链',
    address    text         null comment '地址',
    createTime bigint       null comment '记录创建的时间戳',
    type       tinyint      not null comment '为1时content是推文，为2时content是更新数据',
    constraint twitter_info_pk
        unique (tweetsId)
)
    comment '存储推特相关信息的表' row_format = COMPRESSED;

