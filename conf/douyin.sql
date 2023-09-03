SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

DROP TABLE IF EXISTS `follows`;
CREATE TABLE `follows` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `user_id` bigint(20) NOT NULL COMMENT '用户id',
  `follower_id` bigint(20) NOT NULL COMMENT '关注的用户',
  `cancel` tinyint(4) NOT NULL DEFAULT '0' COMMENT '默认关注为0，取消关注为1',
  PRIMARY KEY (`id`),
  UNIQUE KEY `userIdToFollowerIdIdx` (`user_id`,`follower_id`) USING BTREE,
  KEY `FollowerIdIdx` (`follower_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1096 DEFAULT CHARSET=utf8 COMMENT='关注表';

DROP PROCEDURE IF EXISTS `addFollowRelation`;
delimiter ;;
CREATE PROCEDURE `addFollowRelation`(IN user_id bigint,IN follower_id bigint)
BEGIN
	#Routine body goes here...
	# 声明记录个数变量。
	DECLARE cnt INT DEFAULT 0;
	# 获取记录个数变量。
	SELECT COUNT(1) FROM follows f where f.user_id = user_id AND f.follower_id = follower_id INTO cnt;
	# 判断是否已经存在该记录，并做出相应的插入关系、更新关系动作。
	# 插入操作。
	IF cnt = 0 THEN
		INSERT INTO follows(`user_id`,`follower_id`) VALUES(user_id,follower_id);
	END IF;
	# 更新操作
	IF cnt != 0 THEN
		UPDATE follows f SET f.cancel = 0 WHERE f.user_id = user_id AND f.follower_id = follower_id;
	END IF;
END
;;
delimiter ;

DROP PROCEDURE IF EXISTS `delFollowRelation`;
delimiter ;;
CREATE PROCEDURE `delFollowRelation`(IN `user_id` bigint,IN `follower_id` bigint)
BEGIN
	#Routine body goes here...
	# 定义记录个数变量，记录是否存在此关系，默认没有关系。
	DECLARE cnt INT DEFAULT 0;
	# 查看是否之前有关系。
	SELECT COUNT(1) FROM follows f WHERE f.user_id = user_id AND f.follower_id = follower_id INTO cnt;
	# 有关系，则需要update cancel = 1，使其关系无效。
	IF cnt = 1 THEN
		UPDATE follows f SET f.cancel = 1 WHERE f.user_id = user_id AND f.follower_id = follower_id;
	END IF;
END
;;
delimiter ;

SET FOREIGN_KEY_CHECKS = 1;
