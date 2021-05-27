# Pokemon Team Builder
### Group Members: Toan Vo, Jin Ning Huang, Hongfei Lin, Thomas Gerber

## Project Pitch:
There are many ways to engage with pokemon fans across the world, whether that be through one of the many Nintendo/Gamefreak released pokemon games, or through pokemon based forums and websites such as Pokemon Showdown. In the case of Pokemon Showdown specifically, users cannot save specific teams of pokemon that they want to use in battle. They are forced to rebuild their desired team every time they want to play. This is where the problem comes into play, and our primary goal is to solve it.
Our project revolves around creating an online pokemon database and API. Pokemon fans would be the main target audience because there is a website called pokemonshowdown.com that lets you create teams and battle with them. However, there is no way to store teams online; when you make a team it is stored locally. We want to build this because when a person wants to access a team they already made but on a new computer, they can’t access it. Through our platform, our audience could save multiple teams to their user account, making it easy for them to pick specific teams as soon as they login, eliminating the issue of non-persistent teams in Pokemon Showdown.
Ultimately, we may add even more features to our enhanced Pokemon Showdown user experience, such as statistics about your matches and battles, and your most successful team composition. Yet, the primary focus of our app will be utilizing our API and database to allow users to save, access, and cycle through their favorite teams of pokemon.

## Project Description:
Our project revolves around creating an online pokemon database. Pokemon fans would be the main target audience because there is a website called pokemonshowdown.com that lets you create teams and battle with them. However, there is no way to store teams online; when you make a team it is stored locally. We want to build this because when a person wants to access a team they already made but on a new computer, they can’t access it.

## Technical Description:

<img width="1044" alt="Screen Shot 2021-05-17 at 9 01 33 PM" src="https://user-images.githubusercontent.com/43588965/118588359-13f62880-b753-11eb-8391-084ae221cbe9.png">

Priority | User | Description |
------------ | ------------- | ---------------
P0 | As a user | I want to create pokemon teams that are stored online. A database will store this information. I can access these teams through the **rest API**. <br> **Endpoint**: /teams - looks at teams |
P0 | As a user | I want to create a username with my email and password. My trainerID will be automatically created, stored in **MySQL Database and Redis Store**, and the ID be used for teams. <br> **Endpoint**: /users |
P1 | As a user | Even though we can store the teams, a big part is selecting the movesets for each pokemon. This will take a lot of coding database rules and restrictions, which is why it’s labelled as P1. <br> **Endpoint**: /teams/moves - picks moves for pokemon |
P1 | As a user | Stats are the next thing to tackle and that is its own page under pokemon selection. <br> **Endpoint**: /teams/stats |




![image](https://user-images.githubusercontent.com/49173893/119774721-976cf380-be90-11eb-91f2-23eca04ae9bd.png)
