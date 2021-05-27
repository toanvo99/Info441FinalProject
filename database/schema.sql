create table if not exists Trainer (
    TrainerID int not null auto_increment primary key,
    Username varchar(255) not null unique,
    Email varchar(128) not null unique,
    NumWins int not null
)

create table if not exists Team (
    TeamID int not null auto_increment primary,
    TrainerID int foreign key references Trainer(TrainerID)
)

create table if not exists TeamPokemon (
    TeamID int foreign key references Team(TeamID)
    pokemonSlot1 varchar(128) not null,
    pokemonSlot2 varchar(128) null,
    pokemonSlot3 varchar(128) null,
)

create table if not exists Species (
    SpeciesID int not null primary key,
    Type1ID int not null foreign key references Type(TypeID),
    Type2ID int null foreign key references Type(TypeID),
    SpeciesName varchar(128) not null
)

create table if not exists Type (
    TypeID int primary key,
    TypeName varchar(128)
)

