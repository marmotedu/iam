-- Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
-- Use of this source code is governed by a MIT style
-- license that can be found in the LICENSE file.

CREATE DATABASE  IF NOT EXISTS `iam` /*!40100 DEFAULT CHARACTER SET utf8 */;
USE `iam`;
-- MySQL dump 10.13  Distrib 5.7.29, for Win64 (x86_64)
--
-- Host: 106.52.30.200    Database: iam
-- ------------------------------------------------------
-- Server version	5.5.5-10.5.9-MariaDB

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `policy`
--

DROP TABLE IF EXISTS `policy`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `policy` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `instanceID` varchar(32) DEFAULT NULL,
  `name` varchar(45) NOT NULL,
  `username` varchar(255) NOT NULL,
  `policyShadow` longtext DEFAULT NULL,
  `extendShadow` longtext DEFAULT NULL,
  `createdAt` timestamp NOT NULL DEFAULT current_timestamp(),
  `updatedAt` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  UNIQUE KEY `instanceID_UNIQUE` (`instanceID`),
  KEY `fk_policy_user_idx` (`username`),
  CONSTRAINT `fk_policy_user` FOREIGN KEY (`username`) REFERENCES `user` (`name`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=47 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `policy`
--

LOCK TABLES `policy` WRITE;
/*!40000 ALTER TABLE `policy` DISABLE KEYS */;
/*!40000 ALTER TABLE `policy` ENABLE KEYS */;
UNLOCK TABLES;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_general_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
/*!50003 CREATE*/ /*!50017 DEFINER=`iam`@`127.0.0.1`*/ /*!50003 TRIGGER `iam`.`policy_BEFORE_DELETE` BEFORE DELETE ON `policy` FOR EACH ROW
BEGIN
	insert into policy_audit values(old.id, old.instanceID, old.name, old.username, old.policyShadow, old.extendShadow, old.createdAt, old.updatedAt, curtime());
END */;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;

--
-- Table structure for table `policy_audit`
--

DROP TABLE IF EXISTS `policy_audit`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `policy_audit` (
  `id` bigint(20) unsigned NOT NULL,
  `instanceID` varchar(32) DEFAULT NULL,
  `name` varchar(45) NOT NULL,
  `username` varchar(255) NOT NULL,
  `policyShadow` longtext DEFAULT NULL,
  `extendShadow` longtext DEFAULT NULL,
  `createdAt` timestamp NOT NULL DEFAULT current_timestamp(),
  `updatedAt` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `deletedAt` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  PRIMARY KEY (`id`),
  KEY `fk_policy_user_idx` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `policy_audit`
--

LOCK TABLES `policy_audit` WRITE;
/*!40000 ALTER TABLE `policy_audit` DISABLE KEYS */;
/*!40000 ALTER TABLE `policy_audit` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `secret`
--

DROP TABLE IF EXISTS `secret`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `secret` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `instanceID` varchar(32) DEFAULT NULL,
  `name` varchar(45) NOT NULL,
  `username` varchar(255) NOT NULL,
  `secretID` varchar(36) NOT NULL,
  `secretKey` varchar(255) NOT NULL,
  `expires` int(64) unsigned NOT NULL DEFAULT 1534308590,
  `description` varchar(255) NOT NULL,
  `extendShadow` longtext DEFAULT NULL,
  `createdAt` timestamp NOT NULL DEFAULT current_timestamp(),
  `updatedAt` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  UNIQUE KEY `instanceID_UNIQUE` (`instanceID`),
  KEY `fk_secret_user_idx` (`username`),
  CONSTRAINT `fk_secret_user` FOREIGN KEY (`username`) REFERENCES `user` (`name`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=22 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `secret`
--

LOCK TABLES `secret` WRITE;
/*!40000 ALTER TABLE `secret` DISABLE KEYS */;
/*!40000 ALTER TABLE `secret` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user`
--

DROP TABLE IF EXISTS `user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `user` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `instanceID` varchar(32) DEFAULT NULL,
  `name` varchar(45) NOT NULL,
  `status` int(1) DEFAULT 1 COMMENT '1:可用，0:不可用',
  `nickname` varchar(30) NOT NULL,
  `password` varchar(255) NOT NULL,
  `email` varchar(256) NOT NULL,
  `phone` varchar(20) DEFAULT NULL,
  `isAdmin` tinyint(1) unsigned NOT NULL DEFAULT 0 COMMENT '1: administrator\\\\n0: non-administrator',
  `extendShadow` longtext DEFAULT NULL,
  `loginedAt` timestamp NULL DEFAULT NULL COMMENT 'last login time',
  `createdAt` timestamp NOT NULL DEFAULT current_timestamp(),
  `updatedAt` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_name` (`name`),
  UNIQUE KEY `instanceID_UNIQUE` (`instanceID`)
) ENGINE=InnoDB AUTO_INCREMENT=38 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user`
--

LOCK TABLES `user` WRITE;
/*!40000 ALTER TABLE `user` DISABLE KEYS */;
INSERT INTO `user` VALUES (1,'user-lingfei','admin',1,'admin','$2a$10$WnQD2DCfWVhlGmkQ8pdLkesIGPf9KJB7N1mhSOqulbgN7ZMo44Mv2','admin@foxmail.com','1812884xxxx',1,'{}',now(),'2021-05-27 10:01:40','2021-05-05 21:13:14');
/*!40000 ALTER TABLE `user` ENABLE KEYS */;
UNLOCK TABLES;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_general_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'PIPES_AS_CONCAT,ANSI_QUOTES,ONLY_FULL_GROUP_BY,NO_AUTO_VALUE_ON_ZERO,STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
/*!50003 CREATE*/ /*!50017 DEFINER="iam"@"127.0.0.1"*/ /*!50003 TRIGGER `iam`.`user_BEFORE_DELETE` BEFORE DELETE ON `user` FOR EACH ROW
BEGIN
	delete from secret where username = old.name;
    delete from policy where username = old.name;
END */;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;

--
-- Dumping events for database 'iam'
--

--
-- Dumping routines for database 'iam'
--
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2021-05-27 18:07:55
