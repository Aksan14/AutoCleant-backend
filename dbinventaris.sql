-- MySQL dump 10.13  Distrib 8.0.43, for Linux (x86_64)
--
-- Host: localhost    Database: resetbph_inventaris
-- ------------------------------------------------------
-- Server version	8.0.43-0ubuntu0.22.04.1

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `inventaris`
--

DROP TABLE IF EXISTS `inventaris`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `inventaris` (
  `id` int NOT NULL AUTO_INCREMENT,
  `nama_barang` varchar(100) NOT NULL,
  `kategori` varchar(50) NOT NULL,
  `jumlah` int NOT NULL,
  `satuan` varchar(20) NOT NULL,
  `kondisi` enum('Baik','Rusak Ringan','Rusak Berat','Dimusnahkan') CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT 'Baik',
  `foto` varchar(255) DEFAULT NULL,
  `status` varchar(20) DEFAULT 'tersedia',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=60 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `inventaris`
--

LOCK TABLES `inventaris` WRITE;
/*!40000 ALTER TABLE `inventaris` DISABLE KEYS */;
INSERT INTO `inventaris` VALUES (58,'Baju','Lainnya',10,'Buah','Dimusnahkan','uploads/download.jpeg','tersedia','2025-09-09 06:13:15','2025-09-12 06:30:22'),(59,'tes','Dapur',3,'Pcs','Baik','uploads/Screenshot from 2025-09-11 01-29-59.png','tersedia','2025-09-12 05:35:22','2025-09-12 06:03:46');
/*!40000 ALTER TABLE `inventaris` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `inventaris_check`
--

DROP TABLE IF EXISTS `inventaris_check`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `inventaris_check` (
  `id` int NOT NULL AUTO_INCREMENT,
  `report_id` int unsigned NOT NULL,
  `inventaris_id` int NOT NULL,
  `kondisi` varchar(50) NOT NULL,
  `keterangan` text,
  `tanggal_cek` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_check_report` (`report_id`),
  KEY `idx_check_inventaris` (`inventaris_id`),
  CONSTRAINT `fk_check_inventaris` FOREIGN KEY (`inventaris_id`) REFERENCES `inventaris` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  CONSTRAINT `fk_check_report` FOREIGN KEY (`report_id`) REFERENCES `inventaris_report` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB AUTO_INCREMENT=70 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `inventaris_check`
--

LOCK TABLES `inventaris_check` WRITE;
/*!40000 ALTER TABLE `inventaris_check` DISABLE KEYS */;
INSERT INTO `inventaris_check` VALUES (69,42,58,'rusak','Hilang 1','2025-09-12 02:42:25','2025-09-12 02:42:24','2025-09-12 02:42:24');
/*!40000 ALTER TABLE `inventaris_check` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `inventaris_report`
--

DROP TABLE IF EXISTS `inventaris_report`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `inventaris_report` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `kode_report` varchar(64) NOT NULL,
  `tanggal_report` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `petugas` varchar(100) NOT NULL,
  `status` enum('draft','final') NOT NULL DEFAULT 'draft',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `kode_report` (`kode_report`)
) ENGINE=InnoDB AUTO_INCREMENT=44 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `inventaris_report`
--

LOCK TABLES `inventaris_report` WRITE;
/*!40000 ALTER TABLE `inventaris_report` DISABLE KEYS */;
INSERT INTO `inventaris_report` VALUES (41,'RPT-d09c351e','2025-09-12 02:38:53','Asepp','final','2025-09-12 02:38:52','2025-09-12 02:39:00'),(42,'RPT-b546801a','2025-09-12 02:41:14','Aksan','final','2025-09-12 02:41:13','2025-09-12 02:42:26'),(43,'RPT-e0fc4f16','2025-09-12 02:43:46','Keorganisasian','draft','2025-09-12 02:43:46','2025-09-12 02:43:46');
/*!40000 ALTER TABLE `inventaris_report` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `peminjaman`
--

DROP TABLE IF EXISTS `peminjaman`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `peminjaman` (
  `id` int NOT NULL AUTO_INCREMENT,
  `inventaris_id` int NOT NULL,
  `nama_peminjam` varchar(100) NOT NULL,
  `tgl_pinjam` date NOT NULL,
  `rencana_kembali` date NOT NULL,
  `jumlah` int NOT NULL DEFAULT '1',
  `tgl_kembali` date DEFAULT NULL,
  `kondisi_setelah` varchar(50) DEFAULT NULL,
  `keterangan_kembali` text,
  `status` varchar(20) NOT NULL DEFAULT 'dipinjam',
  `keterangan` text,
  `foto_bukti_kembali` varchar(255) DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `inventaris_id` (`inventaris_id`),
  CONSTRAINT `peminjaman_ibfk_1` FOREIGN KEY (`inventaris_id`) REFERENCES `inventaris` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=67 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `peminjaman`
--

LOCK TABLES `peminjaman` WRITE;
/*!40000 ALTER TABLE `peminjaman` DISABLE KEYS */;
INSERT INTO `peminjaman` VALUES (62,58,'Aseppp','2025-09-10','2025-09-10',1,'2025-09-10','Rusak Ringan','Baik','selesai','tidak','uploads/Screenshot from 2025-09-10 17-50-47.png','2025-09-10 16:27:44','2025-09-10 16:28:36'),(63,58,'Aseppp','2025-09-10','2025-09-17',1,'2025-09-12','Baik','Baikk','selesai','b','uploads/Screenshot from 2025-09-11 01-29-59.png','2025-09-10 16:49:03','2025-09-12 02:36:21'),(66,58,'Aseppp','2025-09-11','2025-09-29',9,'2025-09-12','Baik','Baik','selesai','Baik','uploads/Screenshot from 2025-09-12 14-13-21.png','2025-09-12 02:40:57','2025-09-12 06:30:22');
/*!40000 ALTER TABLE `peminjaman` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `users` (
  `id` varchar(36) NOT NULL,
  `nra` varchar(20) NOT NULL,
  `password` varchar(255) NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `nra` (`nra`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` VALUES ('4c47f561-1712-4d62-9477-420d173ad7dc','1324011','$2a$10$KxfbG2dj/ueZi25McpWG3u5RLfaBQ/0gv9TCNbgaGD3rIOyhm9Fi2','2025-09-09 06:11:47','2025-09-09 06:11:47'),('a7133cf3-ad45-4f31-87eb-648c671c7d4e','1324013','$2a$10$ncGSM1rJ.zBL0MaV4nUYjeatqXR2KKMnNhVyTcuFMU55znxnnfsgG','2025-09-09 06:11:47','2025-09-09 06:11:47'),('bb799d8e-cd84-4afb-889e-bf7a190b3dd0','1324014','$2a$10$EJtowgUMj/agfOD7CPROBunegfoSMBiSHMe45EJem9H.o1HjJ0F2i','2025-09-09 06:11:47','2025-09-09 06:11:47'),('bf0dca96-f20b-47c3-ac1d-b1c699cae908','1324015','$2a$10$miSGdGnVLYFchKq0vQdgd.hr3ld8W5dREGnwCj3NZaYppFqUm1g76','2025-09-09 06:11:47','2025-09-09 06:11:47');
/*!40000 ALTER TABLE `users` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2025-09-12 14:54:26
