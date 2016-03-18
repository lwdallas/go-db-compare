-- phpMyAdmin SQL Dump
-- version 4.4.10
-- http://www.phpmyadmin.net
--
-- Host: localhost:8889
-- Generation Time: Mar 18, 2016 at 01:50 PM
-- Server version: 5.5.42
-- PHP Version: 7.0.0

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET time_zone = "+00:00";

--
-- Database: `testgdbc`
--

-- --------------------------------------------------------

--
-- Table structure for table `first`
--

CREATE TABLE `first` (
  `id` int(11) NOT NULL,
  `first_name` varchar(40) NOT NULL,
  `middle_name` varchar(40) NOT NULL,
  `last_name` varchar(40) NOT NULL,
  `age` int(11) NOT NULL,
  `birthdate` date NOT NULL,
  `description` varchar(40) DEFAULT NULL,
  `more_info` varchar(40) DEFAULT NULL,
  `addr` varchar(40) DEFAULT NULL,
  `city` varchar(40) NOT NULL
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=latin1;

--
-- Dumping data for table `first`
--

INSERT INTO `first` (`id`, `first_name`, `middle_name`, `last_name`, `age`, `birthdate`, `description`, `more_info`, `addr`, `city`) VALUES
(1, 'Bob', 'A', 'Count', 20, '1994-07-17', 'a little info', 'more info', '123 main st', 'cityville'),
(2, 'Carl', 'Ben', 'Gone', 21, '1995-06-20', 'car salesman', NULL, NULL, 'pleasantville'),
(3, 'Sid', 'Timothy', 'Thomas', 30, '1990-03-13', NULL, NULL, NULL, '');

-- --------------------------------------------------------

--
-- Table structure for table `second`
--

CREATE TABLE `second` (
  `id` int(11) NOT NULL,
  `first_name` varchar(40) NOT NULL,
  `middle_name` varchar(40) NOT NULL,
  `last_name` varchar(40) NOT NULL,
  `age` int(11) NOT NULL,
  `birthdate` date NOT NULL,
  `description` varchar(40) DEFAULT NULL,
  `more_info` varchar(40) DEFAULT NULL,
  `addr` varchar(40) DEFAULT NULL,
  `city` varchar(40) NOT NULL
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=latin1;

--
-- Dumping data for table `second`
--

INSERT INTO `second` (`id`, `first_name`, `middle_name`, `last_name`, `age`, `birthdate`, `description`, `more_info`, `addr`, `city`) VALUES
(1, 'Bob', 'Allen', 'Count', 20, '1994-07-17', 'a little info', 'more info', '123 main st', 'cityville'),
(2, 'Carl', 'Ben', 'Gone', 23, '1995-06-20', 'cashier', 'dentist', 'wacky circle', 'pleasantville'),
(3, 'Sid', 'Timothy', 'Thomas', 30, '1990-03-13', NULL, NULL, NULL, ''),
(4, 'Andrew', 'Dishional', 'Guy', 18, '2000-01-03', 'helicopter pilot', 'pink hair', '555 one street', 'onesville');

--
-- Indexes for dumped tables
--

--
-- Indexes for table `first`
--
ALTER TABLE `first`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `second`
--
ALTER TABLE `second`
  ADD PRIMARY KEY (`id`);

--
-- AUTO_INCREMENT for dumped tables
--

--
-- AUTO_INCREMENT for table `first`
--
ALTER TABLE `first`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT,AUTO_INCREMENT=4;
--
-- AUTO_INCREMENT for table `second`
--
ALTER TABLE `second`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT,AUTO_INCREMENT=5;
