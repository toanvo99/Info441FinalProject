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
create table if not exists [Type] (
    TypeID int auto_increment primary key not null,
    TypeName varchar(128) not null
);

/*outer*/
create table if not exists Damage (
    DamagetID int not null auto_increment primary key,
    DamageType varchar(25) not null,
    DamageDesc varchar(500) not null
)

/*outer*/
create table if not exists [Move] (
    MoveID int not null auto_increment primary key,
    TypeID int not null foreign key references [Type](TypeID),
    DamageID int not null foreign key references Damage(DamageID),
    MoveDesc varchar(500) not null,
    [Power] int not null
);

/*stat*/
create table if not exists [Stat] (
    StatID int not null auto_increment primary key,
    StatName varchar(25) not null,
);

/*type*/
create table if not exists Species (
    SpeciesID int not null auto_increment primary key,
    Type1ID int not null foreign key references Type(TypeID),
    Type2ID int null foreign key references Type(TypeID),
    SpeciesName varchar(128) not null
);

/*species*/
create table if not exists Pokemon (
    PokemonID int not null auto_increment primary key,
    SpeciesID int foreign references Species(SpeciesID) not null
);

/*trainer*/
create table if not exists Team (
    TeamID int not null auto_increment primary key,
    TrainerID int foreign key references Trainer(TrainerID)
);


create table if not exists MoveSet (
    MoveSetID int not null auto_increment primary key
);

create table if not exists MoveSetMove (
    MoveSetID int not null foreign key references MoveSet(MoveSetID),
    MoveID int not null foreign key references [Move](MoveID)
);



create table if not exists StatSet (
    StatSetID int not null auto_increment primary key
);

create table if not exists StatSetStat (
    StatSetID int not null foreign key references StatSet(StatSetID),
    StatID int not null foreign key references [Stat](StatID),
    StatValue int not null
);


create table if not exists PokemonStatSet (
    PokemonID int not null foreign key references Pokemon(PokemonID)
    MoveSetID int not null foreign key references StatSet(StatSetID)
);


create table if not exists PokemonMoveSet (
    PokemonID int not null foreign key references Pokemon(PokemonID)
    MoveSetID int not null foreign key references MoveSet(StatSetID)
);

create table if not exists PokemonTeam (
    PokemonID int not null foreign key references Pokemon(PokemonID),
    MoveSetID int not null foreign key references Team(TeamID)
);



