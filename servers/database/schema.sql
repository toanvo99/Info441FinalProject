/*outer*/
create table if not exists User (
    TrainerID int not null auto_increment primary key,
    Email varchar(128) not null unique,
    PassHash binary(64) not null,
    UserName varchar(255) not null unique,
    FirstName varchar(128) not null,
    LastName varchar(128) not null,
    PhotoURL varchar(128) not null
);

/*outer*/
create table if not exists `Type` (
    TypeID int auto_increment primary key not null,
    TypeName varchar(128) not null
);

/*outer*/
create table if not exists Damage (
    DamagetID int not null auto_increment primary key,
    DamageType varchar(25) not null,
    DamageDesc varchar(500) not null
);

/*outer*/
create table if not exists `Move` (
    MoveID int not null auto_increment primary key,
    TypeID int not null,
    DamageID int not null,
    MoveName varchar(25) not null,      
    MoveDesc varchar(500) not null,
    `Power` int not null
);

/*stat*/
create table if not exists `Stat` (
    StatID int not null auto_increment primary key,
    StatName varchar(25) not null
);

/*type*/
create table if not exists Species (
    SpeciesID int not null auto_increment primary key,
    Type1ID int not null,
    Type2ID int not null,
    SpeciesName varchar(128) not null,
    Sprite varchar(1000) not null
);

/*species*/
create table if not exists Pokemon (
    PokemonID int not null auto_increment primary key,
    SpeciesID int not null
);

/*trainer*/
create table if not exists Team (
    TeamID int not null auto_increment primary key,
    TrainerID int 
);


create table if not exists MoveSet (
    MoveSetID int not null auto_increment primary key
);

create table if not exists MoveSetMove (
    MoveSetID int not null,
    MoveID int not null 
);



create table if not exists StatSet (
    StatSetID int not null auto_increment primary key
);

create table if not exists StatSetStat (
    StatSetID int not null,
    StatID int not null,
    StatValue int not null
);


create table if not exists PokemonStatSet (
    PokemonID int not null,
    MoveSetID int not null 
);


create table if not exists PokemonMoveSet (
    PokemonID int not null,
    MoveSetID int not null 
);

create table if not exists PokemonTeam (
    PokemonID int not null,
    MoveSetID int not null
);

INSERT INTO `Type` (TypeName) VALUES ('Fire'), ('Water');

