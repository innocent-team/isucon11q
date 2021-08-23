
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
