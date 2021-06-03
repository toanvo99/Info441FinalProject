import React from 'react';
import { useState } from 'react';

export function PokemonForms(props) {
    let user = props.user;
    const [pokemon, setPokemon] = useState("");
    const [sprites, setSprites] = useState([]);
    const [submitted, setSubmitted] = useState(false);
    user = "hey there";
    const submitTeam = (event) => {
        event.preventDefault();
        console.log("succesful submit! ");
        setSprites([]);
        setSubmitted(true);
    }
    if (user) {
        return (
            <div className="poke-form">
                <PokemonSearchForm 
                    submitTeam={submitTeam} 
                    pokemon={pokemon}  
                    setPokemon={setPokemon} 
                    sprites={sprites} 
                    setSprites={setSprites}/>
                {submitted ? 
                    <h3 id="success_submit" >
                        Succesfully Submitted!
                    </h3>
                    :
                    <div>  
                    </div>
                }
            </div>
        )
    } else {
        return (
            <div className="poke-prompt">
                <h3>
                    Please sign in to create a new team!
                </h3>
            </div>
        )
    }
}

function PokemonSearchForm(props) {
    const handlePokeSearch= (event) => {
        let newValue = event.target.value
        props.setPokemon(newValue);
    }
    return (
        /*<form>
            <label htmlFor="change_description">Search Pokemon Database:</label>
            <input id="search_pokemon" className="form-control" type="text" name="term" aria-label="Pokemon Search" placeholder="Ex: piplup" value={props.pokemon} onChange={handlePokeSearch}></input>
        </form>*/
        <form action="/" method="get" onSubmit={props.submitTeam}>
            <label htmlFor="header-search">
                <span className="visually-hidden">Search pokemon</span>
            </label>
            <input
                type="text"
                id="header-search"
                placeholder="Ex: piplup"
                value={props.pokemon}
                onChange={handlePokeSearch}
                name="s" 
            />
            <PokemonSearchButton  pokemon={props.pokemon} 
                sprites={props.sprites}
                setIsTeam={props.isTeam}
                setSprites={props.setSprites}/>
            {props.sprites.length !== 0 ?
                <input id="submit" onSubmit={props.submitTeam} type="submit" value="Submit Team" />
                :
                <div></div>
            }
        </form>
    )
}

function PokemonSearchButton(props) {
    const handleClickSearch = (pokemon) => {
        // Here, we need to make an api call to
        // retrieve the pokemon. If it fails,
        // do not add the pokemon to our team.

        // TESTING CODE/ Code for demo.
        let paths = props.sprites;
        if (pokemon && props.sprites.length < 6) {
            //print(pokemon);
            if (pokemon === "empoleon") {
                paths.push("./img/emp_sprite.png");
                props.setSprites(paths);
            } else if (pokemon === "heatran") {
                paths.push("./img/heatran_sprite.png");
                props.setSprites(paths);
            } else if (pokemon === "lugia") {
                paths.push("./img/lugia_sprite.png");
                props.setSprites(paths);
            } else if (pokemon === "regigigas") {   
                paths.push("./img/regi_sprite.png");
                props.setSprites(paths);
            } else if (pokemon === "slowking") {
                paths.push("./img/slowking_sprite.png");
                props.setSprites(paths);
            } else if (pokemon === "steelix") {
                paths.push("./img/steelix_sprite.png");
                props.setSprites(paths);
            }
        } else {

        }
    }
    const handleClick = (event) => {
        event.preventDefault();
        handleClickSearch(props.pokemon)
      }
      
    return (
      <div>
        <button id="search_button" type="submit" className="btn btn-primary" onClick={handleClick}>Search!</button>
        <div>
            {props.sprites.map(item => (
                <img src={item} alt="poke sprite" />
            ))} 
          </div>
      </div>
    )
  }
