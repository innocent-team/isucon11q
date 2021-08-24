CREATE TABLE `neo_isu_condition` (
  `jia_isu_uuid` CHAR(36) NOT NULL,
  `timestamp` DATETIME NOT NULL,
  `is_sitting` TINYINT(1) NOT NULL,
  `condition` VARCHAR(255) NOT NULL,
  `message` VARCHAR(255) NOT NULL,
  `created_at` DATETIME(6) DEFAULT CURRENT_TIMESTAMP(6),

  PRIMARY KEY (`jia_isu_uuid`, `timestamp`)
) ENGINE=InnoDB DEFAULT CHARACTER SET=utf8mb4;

INSERT INTO `neo_isu_condition` SELECT `jia_isu_uuid`, `timestamp`, `is_sitting`, `condition`, `message`, `created_at` FROM `isu_condition`;
DROP TABLE `isu_condition`;
RENAME TABLE `neo_isu_condition` TO `isu_condition`;

ALTER TABLE `isu_condition` ADD (
    `is_broken` TINYINT(1) AS (`condition` REGEXP 'is_broken=true') STORED,
    `is_dirty` TINYINT(1) AS (`condition` REGEXP 'is_dirty=true') STORED,
    `is_overweight` TINYINT(1) AS (`condition` REGEXP 'is_overweight=true') STORED,
    `condition_level` TINYINT(1) AS (`is_dirty` + `is_overweight` + `is_broken`) STORED
);

ALTER TABLE `isu` ADD COLUMN  `last_condition_timestamp` DATETIME NOT NULL DEFAULT '2000-01-01 00:00:00';
