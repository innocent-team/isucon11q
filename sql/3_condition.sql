ALTER TABLE `isu_condition` ADD (
    `is_broken` TINYINT(1) AS (`condition` REGEXP 'is_broken=true') STORED,
    `is_dirty` TINYINT(1) AS (`condition` REGEXP 'is_dirty=true') STORED,
    `is_overweight` TINYINT(1) AS (`condition` REGEXP 'is_overweight=true') STORED,
    `condition_level` TINYINT(1) AS (`is_dirty` + `is_overweight` + `is_broken`) STORED
);
