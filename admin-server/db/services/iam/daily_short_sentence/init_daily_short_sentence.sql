-- iam/daily_short_sentence 初始化数据
-- 从 admin-server/db/data.sql 拆分而来，见 docs/15-service-boundaries.md 第 4 节

-- 9.1 每日短句模块初始化（菜单、权限、接口）
SET @parent_menu_id = COALESCE(
  (SELECT `id` FROM `admin_menu` WHERE `path` = '/admin/temp' AND `deleted_at` = 0 LIMIT 1),
  (SELECT `id` FROM `admin_menu` WHERE `id` = 9 AND `deleted_at` = 0 LIMIT 1)
);

INSERT INTO `admin_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES (
    @parent_menu_id,
    '每日短句',
    '/admin/temp/daily_short_sentence',
    'misc/DailyShortSentenceList',
    'ele-Document',
    2,
    0,
    1,
    1,
    UNIX_TIMESTAMP(),
    UNIX_TIMESTAMP(),
    0
) ON DUPLICATE KEY UPDATE `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;

SET @main_menu_id = (SELECT `id` FROM `admin_menu` WHERE `path` = '/admin/temp/daily_short_sentence' AND `deleted_at` = 0 LIMIT 1);

INSERT INTO `admin_menu` (`parent_id`, `name`, `path`, `component`, `icon`, `type`, `order_num`, `visible`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (@main_menu_id, '每日短句 新增按钮', '', '', '', 3, 1, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@main_menu_id, '每日短句 编辑按钮', '', '', '', 3, 2, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (@main_menu_id, '每日短句 删除按钮', '', '', '', 3, 3, 0, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;

SET @create_button_id = (SELECT `id` FROM `admin_menu` WHERE `parent_id`=@main_menu_id AND `name`='每日短句 新增按钮' AND `deleted_at`=0 LIMIT 1);
SET @update_button_id = (SELECT `id` FROM `admin_menu` WHERE `parent_id`=@main_menu_id AND `name`='每日短句 编辑按钮' AND `deleted_at`=0 LIMIT 1);
SET @delete_button_id = (SELECT `id` FROM `admin_menu` WHERE `parent_id`=@main_menu_id AND `name`='每日短句 删除按钮' AND `deleted_at`=0 LIMIT 1);

INSERT INTO `admin_permission` (`name`, `code`, `description`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  ('每日短句列表', 'daily_short_sentence:list', '查看每日短句列表', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('每日短句新增', 'daily_short_sentence:create', '新增每日短句', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('每日短句编辑', 'daily_short_sentence:update', '编辑每日短句', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('每日短句删除', 'daily_short_sentence:delete', '删除每日短句', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;

SET @list_permission_id = (SELECT `id` FROM `admin_permission` WHERE `code`='daily_short_sentence:list' AND `deleted_at`=0 LIMIT 1);
SET @create_permission_id = (SELECT `id` FROM `admin_permission` WHERE `code`='daily_short_sentence:create' AND `deleted_at`=0 LIMIT 1);
SET @update_permission_id = (SELECT `id` FROM `admin_permission` WHERE `code`='daily_short_sentence:update' AND `deleted_at`=0 LIMIT 1);
SET @delete_permission_id = (SELECT `id` FROM `admin_permission` WHERE `code`='daily_short_sentence:delete' AND `deleted_at`=0 LIMIT 1);

INSERT INTO `admin_api` (`name`, `method`, `path`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  ('每日短句列表', 'GET', '/api/v1/daily-short-sentences', '获取每日短句列表', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('每日短句新增', 'POST', '/api/v1/daily-short-sentences', '新增每日短句', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('每日短句编辑', 'PUT', '/api/v1/daily-short-sentences', '编辑每日短句', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ('每日短句删除', 'DELETE', '/api/v1/daily-short-sentences', '删除每日短句', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE `updated_at`=UNIX_TIMESTAMP(), `deleted_at`=0;

SET @list_api_id = (SELECT `id` FROM `admin_api` WHERE `path`='/api/v1/daily-short-sentences' AND `method`='GET' AND `deleted_at`=0 LIMIT 1);
SET @create_api_id = (SELECT `id` FROM `admin_api` WHERE `path`='/api/v1/daily-short-sentences' AND `method`='POST' AND `deleted_at`=0 LIMIT 1);
SET @update_api_id = (SELECT `id` FROM `admin_api` WHERE `path`='/api/v1/daily-short-sentences' AND `method`='PUT' AND `deleted_at`=0 LIMIT 1);
SET @delete_api_id = (SELECT `id` FROM `admin_api` WHERE `path`='/api/v1/daily-short-sentences' AND `method`='DELETE' AND `deleted_at`=0 LIMIT 1);

INSERT INTO `admin_permission_menu` (`permission_id`, `menu_id`, `created_at`, `updated_at`)
VALUES
  (@list_permission_id, @main_menu_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@create_permission_id, @create_button_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@update_permission_id, @update_button_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@delete_permission_id, @delete_button_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())
ON DUPLICATE KEY UPDATE `updated_at`=UNIX_TIMESTAMP();

INSERT INTO `admin_permission_api` (`permission_id`, `api_id`, `created_at`, `updated_at`)
VALUES
  (@list_permission_id, @list_api_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@create_permission_id, @create_api_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@update_permission_id, @update_api_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
  (@delete_permission_id, @delete_api_id, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())
ON DUPLICATE KEY UPDATE `updated_at`=UNIX_TIMESTAMP();

-- 9.2 每日短句初始化数据
INSERT INTO `daily_short_sentence` (`id`, `type`, `content`, `img`, `literature_author`, `convert_img`, `created_at`, `updated_at`, `deleted_at`)
VALUES
  (1, 2, '我愿做一颗永不生锈的螺丝钉。', 'https://t10.baidu.com/it/u=2607170580,4056796944&fm=58', '雷锋', NULL, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (2, 2, '母爱胜于万爱。', 'https://t10.baidu.com/it/u=3325666458,3828073077&fm=58', '莎士比亚', NULL, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (3, 2, 'If you shed tears when you miss the sun,you also miss the stars. 如果你因为失去太阳而落泪，那么你也将失去群星。', 'https://t12.baidu.com/it/u=3583036367,4054301455&fm=58', '泰戈尔', NULL, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (4, 2, '没有在深夜痛哭过的人，不足以谈人生。', 'https://t10.baidu.com/it/u=3894103860,4159876305&fm=58', '托马斯·卡莱尔', NULL, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (5, 2, '那一世，你为蝴蝶，我为落花，花心已碎，蝶翼天涯，那一世，你为繁星，我为月牙，形影相错，空负年华，那一世，你为歌女，我为琵琶，乱世笙歌，深情天下，金戈铁马，水月镜花，容华一刹那，那缕传世的青烟，点缀着你我结缘的童话。不问贵贱，不顾浮华，三千华发，一生牵挂。', 'https://t10.baidu.com/it/u=2066699184,3350866713&fm=58', '佚名', NULL, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (6, 2, '鲜花打扮不出美丽的春天，一个人先进总是单枪匹马，众人先进才能移山填海。', 'https://t11.baidu.com/it/u=3681108379,232235240&fm=58', '雷锋', NULL, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (7, 2, '劳动是一切知识的源泉。', 'https://t11.baidu.com/it/u=3277121875,3598099677&fm=58', '陶铸', NULL, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (8, 2, '一朵鲜花打扮不出美丽的春天，一个人先进总是单一槍一匹马，众人先进才能移山填海。', 'https://t11.baidu.com/it/u=1167556189,1973262450&fm=58', '雷锋', NULL, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (9, 2, '强悍太久，让我软弱很难。', 'https://t11.baidu.com/it/u=3439716001,3622666293&fm=58', '赵丽颖', NULL, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (10, 2, '信仰是石，擦起星星之火；信仰是火，点亮希望之灯；信仰是灯，照亮夜行的路；信仰是路，引你走向黎明。', 'https://t12.baidu.com/it/u=546752963,682448945&fm=58', '佚名', NULL, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (11, 2, '站立着的心，只有努力努力再努力。', 'https://t11.baidu.com/it/u=3748327216,3933540256&fm=58', '张艺兴', NULL, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (12, 2, '读书使人得到一种优雅和风味，这就是读书的整个目的，而只有抱着这种目的的读书才可以叫做艺术，一人读书的目的并不是要"改进心智"，因为当他开始想要改进心智的时候，一切读书的乐趣便丧失净尽了。', 'https://t11.baidu.com/it/u=4257120178,225048982&fm=58', '林语堂', NULL, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (13, 2, '去见你想见的人吧。趁阳光正好，趁微风不噪，趁繁花还未开至荼蘼，趁现在还年轻，还可以走很长很长的路，还能诉说很深很深的思念，趁世界还不那么拥挤，趁飞机还没有起飞，趁现在自己的双手还能拥抱彼此，趁我们还有呼吸。', 'https://t11.baidu.com/it/u=1093319152,1457586368&fm=58', '佚名', NULL, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (14, 2, '孤单不是与生俱来，而是由你爱上一个人的那一刻开始。', 'https://t12.baidu.com/it/u=3224252345,3273330929&fm=58', '张小娴', NULL, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (15, 2, '我谈过最长的恋爱，就是自恋，我爱自己，没有情敌。', 'https://t12.baidu.com/it/u=289577954,149653588&fm=58', '安东尼', NULL, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (16, 2, '上帝创造了整数，所有其余的数都是人造的。', 'https://t11.baidu.com/it/u=4055619370,68771096&fm=58', 'L·克隆内克', NULL, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (17, 2, '帝王：待我君临天下，许你四海为家；国臣：待我了无牵挂，许你浪迹天涯；将军：待我半生戎马，许你共话桑麻；书生：待我功成名达，许你花前月下；侠客：待我名满华夏，许你放歌纵马；琴师：待我弦断音垮，许你青丝白发；面首：待我不再有她，许你淡饭粗茶；情郎：待我高头大马，许你嫁衣红霞；农夫：待我富贵荣华，许你十里桃花；僧人：待我一袭袈裟，许你相思放下。', 'https://t10.baidu.com/it/u=2524492674,127412282&fm=58', '佚名', NULL, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (18, 2, '我是广大劳苦大众当中的一员，我能帮忙人民克服一点困难，是最幸福的。', 'https://t11.baidu.com/it/u=3991612066,1002782314&fm=58', '雷锋', NULL, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (19, 2, '雨下给富人，也下给穷人，下给义人，也下给不义的人；其实，雨并不公道，因为下落在一个没有公道的世界上。', 'https://t10.baidu.com/it/u=578245077,202302703&fm=58', '老舍', NULL, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (20, 2, '只要能培一朵花，就不妨做做会朽的腐草。', 'https://t12.baidu.com/it/u=472244129,624373840&fm=58', '鲁迅', NULL, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (21, 2, '当你做对的时候，没有人会记得；当你做错的时候，连呼吸都是错。', 'https://t12.baidu.com/it/u=3520060508,2659517975&fm=58', '郭敬明', NULL, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (22, 2, '故事的开头总是这样，适逢其会，猝不及防。故事的结局总是这样，花开两朵，天各一方。', 'https://t11.baidu.com/it/u=4192581556,724875387&fm=58', '张嘉佳', NULL, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (23, 2, '下定决心，不怕牺牲，排除万难，去争取胜利。', 'https://t12.baidu.com/it/u=1394716907,535750628&fm=58', '毛泽东', NULL, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (24, 2, '青春是明知道错了，偏要任性到底！', 'https://t12.baidu.com/it/u=3520060508,2659517975&fm=58', '何炅', NULL, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (25, 2, '坚持自己做的事情就可以了，时间会告诉你你的选择正确与否。', 'https://t10.baidu.com/it/u=4208160840,613911946&fm=58', '金星', NULL, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (26, 2, '但愿每次回忆，对 生活都不认为负疚。', 'https://t12.baidu.com/it/u=3224252345,3273330929&fm=58', '', NULL, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (27, 2, '宁肯少些，但要好些。', 'https://t12.baidu.com/it/u=3208410425,757503539&fm=58', '列宁', NULL, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (28, 2, '也许路途很遥远，也许这条路很危险，但是我眼中的风景，是你想像不到的耀眼。', 'https://t12.baidu.com/it/u=472244129,624373840&fm=58', '杨幂', NULL, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (29, 2, '别低头，王冠会掉，别流泪，坏人会笑。', 'https://t12.baidu.com/it/u=3687046766,3009085819&fm=58', '佚名', NULL, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (30, 2, '如果你是一滴水，你是否滋润了一寸土地？如果你是一线阳光，你是否照亮了一分黑暗？如果你是一粒粮食，你是否哺育了有用的生命？如果你是最小的一颗螺丝钉，你是否永远坚守你生活的岗位。', 'https://t11.baidu.com/it/u=3993765373,4203836679&fm=58', '雷锋', NULL, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (31, 2, '我们记得，马吕斯便是从这儿开始的，狂热的恋情忽然出现，并把他推到了种种无目的和无基础的幻想中，他出门仅仅为了去胡思乱想，缓慢的渍染，喧闹而淤止的深渊，并且，随着工作的减少，需要增加了，这是一条规律，处于梦想状态中的人自然是不节约、不振作的，弛懈的精神经受不住紧张的生活，在这种生活方式中，有坏处也有好处，因为慵懒固然有害，慷慨却是健康和善良的，但是不工作的人，穷而慷慨高尚，那是不可救药的，财源涸竭，费用急增， 这是一条导向绝境的下坡路，在这方面，最诚实和最稳定的人也能跟最软弱和最邪恶的人一样往下滑，一直滑到两个深坑中的一个里去：自杀或是犯罪。', 'https://t12.baidu.com/it/u=1579004585,3659949784&fm=58', '雨果', NULL, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (32, 2, '阅读的最大理由是想摆脱平庸，早一天就多一份人生的精彩；迟一天就多一天平庸的困扰。', 'https://t10.baidu.com/it/u=174649884,718480879&fm=58', '余秋雨', NULL, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE
  `type`=VALUES(`type`),
  `content`=VALUES(`content`),
  `img`=VALUES(`img`),
  `literature_author`=VALUES(`literature_author`),
  `convert_img`=VALUES(`convert_img`),
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

-- 9.3 飞书用户角色默认权限：daily_short_sentence 是 init-dev-db.sh 里 iam 域按依赖顺序
-- 排列的最后一个模块，只有到这里 daily_short_sentence:list（本文件）和 file:list
-- （iam/permission）才都已经存在，飞书角色（iam/role）赋权只能放在这里，不能放在
-- role 模块自己的初始化文件里。
SET @feishu_role_id = (SELECT `id` FROM `admin_role` WHERE `code` = 'feishu' AND `deleted_at` = 0 LIMIT 1);

INSERT INTO `admin_role_permission` (`role_id`, `permission_id`, `created_at`, `updated_at`)
SELECT @feishu_role_id, p.`id`, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()
FROM `admin_permission` p
WHERE p.`code` IN ('daily_short_sentence:list', 'file:list')
  AND NOT EXISTS (
    SELECT 1 FROM `admin_role_permission` rp
    WHERE rp.`role_id` = @feishu_role_id AND rp.`permission_id` = p.`id`
  );

