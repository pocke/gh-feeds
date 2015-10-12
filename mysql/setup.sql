-- MySQL Workbench Forward Engineering

SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0;
SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0;
SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='TRADITIONAL,ALLOW_INVALID_DATES';

-- -----------------------------------------------------
-- Schema ghfeeds
-- -----------------------------------------------------

-- -----------------------------------------------------
-- Schema ghfeeds
-- -----------------------------------------------------
CREATE SCHEMA IF NOT EXISTS `ghfeeds` DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci ;
USE `ghfeeds` ;

-- -----------------------------------------------------
-- Table `ghfeeds`.`users`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `ghfeeds`.`users` (
  `id` INT NOT NULL COMMENT 'GitHub ID',
  `name` VARCHAR(45) NOT NULL COMMENT 'User name of GitHub.\nthis is \'login\' on result of GitHub API.',
  `auth` VARCHAR(45) NOT NULL COMMENT '',
  PRIMARY KEY (`id`)  COMMENT '',
  UNIQUE INDEX `id_UNIQUE` (`id` ASC)  COMMENT '',
  UNIQUE INDEX `username_UNIQUE` (`name` ASC)  COMMENT '')
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `ghfeeds`.`feeds`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `ghfeeds`.`feeds` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT '',
  `type` VARCHAR(45) NOT NULL COMMENT 'WatchEvent, IssueCommentEvent, ...',
  `published_at` DATETIME(5) NOT NULL COMMENT '',
  `user_id` INT NOT NULL COMMENT '',
  `html` TEXT NOT NULL COMMENT '',
  `author_name` VARCHAR(45) NOT NULL COMMENT '',
  `url` TEXT NOT NULL COMMENT '',
  `image_url` TEXT NOT NULL COMMENT '',
  PRIMARY KEY (`id`, `user_id`)  COMMENT '',
  INDEX `fk_feed_entries_users_idx` (`user_id` ASC)  COMMENT '',
  CONSTRAINT `fk_feed_entries_users`
    FOREIGN KEY (`user_id`)
    REFERENCES `ghfeeds`.`users` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB;


SET SQL_MODE=@OLD_SQL_MODE;
SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS;
SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS;
