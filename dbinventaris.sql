-- MySQL dump 10.13  Distrib 8.0.19, for Win64 (x86_64)
--
-- Host: localhost    Database: resetbph_inventaris
-- ------------------------------------------------------
-- Server version	8.0.43

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
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=58 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `inventaris`
--

LOCK TABLES `inventaris` WRITE;
/*!40000 ALTER TABLE `inventaris` DISABLE KEYS */;
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
  PRIMARY KEY (`id`),
  KEY `idx_check_report` (`report_id`),
  KEY `idx_check_inventaris` (`inventaris_id`),
  CONSTRAINT `fk_check_inventaris` FOREIGN KEY (`inventaris_id`) REFERENCES `inventaris` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  CONSTRAINT `fk_check_report` FOREIGN KEY (`report_id`) REFERENCES `inventaris_report` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB AUTO_INCREMENT=62 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `inventaris_check`
--

LOCK TABLES `inventaris_check` WRITE;
/*!40000 ALTER TABLE `inventaris_check` DISABLE KEYS */;
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
  PRIMARY KEY (`id`),
  UNIQUE KEY `kode_report` (`kode_report`)
) ENGINE=InnoDB AUTO_INCREMENT=12 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `inventaris_report`
--

LOCK TABLES `inventaris_report` WRITE;
/*!40000 ALTER TABLE `inventaris_report` DISABLE KEYS */;
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
  `tgl_kembali` date DEFAULT NULL,
  `kondisi_setelah` varchar(50) DEFAULT NULL,
  `status` varchar(20) NOT NULL DEFAULT 'dipinjam',
  `keterangan` text,
  `foto_bukti_kembali` varchar(255) DEFAULT NULL,
  `keterangan_kembali` text DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `inventaris_id` (`inventaris_id`),
  CONSTRAINT `peminjaman_ibfk_1` FOREIGN KEY (`inventaris_id`) REFERENCES `inventaris` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=26 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `peminjaman`
--

LOCK TABLES `peminjaman` WRITE;
/*!40000 ALTER TABLE `peminjaman` DISABLE KEYS */;
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
  PRIMARY KEY (`id`),
  UNIQUE KEY `nra` (`nra`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` VALUES ('4c47f561-1712-4d62-9477-420d173ad7dc','1324011','$2a$10$KxfbG2dj/ueZi25McpWG3u5RLfaBQ/0gv9TCNbgaGD3rIOyhm9Fi2'),('a7133cf3-ad45-4f31-87eb-648c671c7d4e','1324013','$2a$10$ncGSM1rJ.zBL0MaV4nUYjeatqXR2KKMnNhVyTcuFMU55znxnnfsgG'),('bb799d8e-cd84-4afb-889e-bf7a190b3dd0','1324014','$2a$10$EJtowgUMj/agfOD7CPROBunegfoSMBiSHMe45EJem9H.o1HjJ0F2i'),('bf0dca96-f20b-47c3-ac1d-b1c699cae908','1324015','$2a$10$miSGdGnVLYFchKq0vQdgd.hr3ld8W5dREGnwCj3NZaYppFqUm1g76');
/*!40000 ALTER TABLE `users` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping routines for database 'resetbph_inventaris'
--
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2025-09-03 20:00:47
