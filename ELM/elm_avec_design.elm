

module Jeu exposing(..)

import Browser
import Html exposing (..)
import Html.Attributes exposing (..)
import Html exposing (Html, Attribute, div, input, text, button)
import Html.Events exposing (..)
import Http
import Random 



import Json.Decode exposing (..)


-- MODEL

type State = Failure String | Loading | Success String             --création du type State qui peut prendre 3 valeurs

type alias Datas = { word : String, meanings : List Meaning}              --création d'alias de type : Datas, Meaning et Definition
type alias Meaning = {partOfSpeech : String, definitions : List Definition}    
type alias Definition = {definition : String}

type alias Model =             -- création du modèle utilisé par l'application
    { http : State
    , content : String 
    , lWords : List String
    , word : String
    , jSon : State
    , datas : List Datas
    }

init : () -> (Model, Cmd Msg)             
init _ =                                              -- initialisation du modèle 
  ( Model Loading "" [] "" Loading []
  , Http.get
      { url = " http://localhost:8000/monFichier.txt"                   -- récupération du fichier txt contenant les mots 
      , expect = Http.expectString Words
      }
  )

-- UPDATE


type Msg                                      -- définition du type Msg pour les msg envoyés à l'application
  = Words (Result Http.Error String)
  | RandomWord Int                                -- Msg peut prendre différentes valeurs : Words,RandomWord,GotJson et Change 
  | GotJson (Result Http.Error (List Datas))
  | Change String
  | Next 
  
update : Msg -> Model -> (Model, Cmd Msg)               -- update prends un msg et un model en argument
update msg model =
  case msg of                             -- avec case on traite en fonction du message recu 
    Words result ->
      case result of
        Ok words ->
          ({ model | lWords = String.words words , http = Success "" } , randomWord model)      --  si la réponse http est un succès on met à jour http 

        Err error ->
          ({model | http = Failure (toString error)}, Cmd.none)         -- si il y a une erreur alors on met à jour l'état de http à échec
          
    RandomWord index -> case (getElemList model.lWords index) of                  -- récupère un mot dans la liste à l'index spécifié par le msg
                                Nothing -> (model, randomWord model)            -- si le mot n'existe pas on refait appel à randomword
                                Just x -> ({ model | word = x }, Http.get {url = ("https://api.dictionaryapi.dev/api/v2/entries/en/" ++ x)  , expect = Http.expectJson GotJson listDatasDecoder})
                                                      -- si le mot existe on met à jour la valeur Word du modèle puis envois une requete pour récuperer les définitions 
    GotJson result -> case result of
                            Ok data-> ({ model | jSon = Success "" , datas = data} , Cmd.none)  -- si il y a succès alors on met à jour datas dans le modèle et met à jour l'état de jSon.

                            Err error -> ({ model | jSon = Failure (toString error) } , randomWord model)   -- si il y a une erreur on met à jour l'état de jSon et on rappel randomWorld pour avoir un autre mot.

    Change newContent -> ({model | content = newContent}, Cmd.none)   -- on met à jour la valeur de content du modèle.
    
    Next  ->( Model Loading "" [] "" Loading [], Http.get { url = " http://localhost:8000/monFichier.txt", expect = Http.expectString Words })
      
randomWord : Model -> Cmd Msg
randomWord model = Random.generate RandomWord (Random.int 1 (1000))        -- on utilise randomWorld pour générer un msg RandomWorld avec l'index entre 1 et 1000.


-- VIEW
viewWord : Model -> Html Msg   -- prend un model en entrée
viewWord model = 

  case model.http of                -- on vérifie l'état de http 
    Loading -> text "Chargement de la page"          -- en fonction de l'état on affiche des messages différents.
    
    Failure error -> text ("Chargement impossible : " ++ error) 

    Success good -> overlay model (                  -- en cas de succès on fait appel à la fonction overlay définie plus bas
      case model.jSon of
        Success veryGood -> [div [] (List.map viewData model.datas)]    -- on affiche les données dans une balise div
        Loading -> [text "Chargement"]                  -- on affiche chargement
        Failure error -> [text ("échec " ++ error)] )    -- on affiche échec et l'erreur.

viewData : Datas -> Html Msg      -- prend Datas en entrée
viewData data =
    div []
        [                                    -- on créer une balise vide où on ajoute la liste ul qui est remplie avec list map et viewmeaning à partir des Datas d'entrée.
         ul [] (List.map viewMeaning data.meanings)
        ]

viewMeaning : Meaning -> Html Msg   -- prend  un objet meaning en entrée
viewMeaning meaning =
    div []                            -- créer une balise vide  et y ajoute un en-tête qui contient partOfSpeech de meaning
        [ h4 [] [text meaning.partOfSpeech]
        , ul [] (List.map viewDefinition meaning.definitions)    -- ajout de la liste ul qui est remplie avec list map et viewdefinition 
        ]

viewDefinition : Definition -> Html Msg    -- prend un objet definitin en entrée 
viewDefinition definition =             
    li [] [text definition.definition]        -- créer une balise vide li et y ajoute le texte de la définition.





-- SUBSCRIPTIONS


subscriptions : Model -> Sub Msg
subscriptions model =
  Sub.none




          
-- JSON                     définition des fonctions listDatasDecoder, definitionDecoder, datasDecoder et meaningDecoder

listDatasDecoder : Decoder (List Datas)
listDatasDecoder = Json.Decode.list datasDecoder
definitionDecoder : Decoder Definition
definitionDecoder = Json.Decode.map Definition (field "definition" string)
datasDecoder : Decoder Datas
datasDecoder = map2 Datas (field "word" string)(field "meanings" <| Json.Decode.list meaningDecoder)
meaningDecoder : Decoder Meaning
meaningDecoder = map2 Meaning (field "partOfSpeech" string)(field "definitions" <| Json.Decode.list definitionDecoder)


-- HELPERS                                            définition des fonctions getElemList, toString,textDatas,textDef,textMeaning et overlay
getElemList : List a -> Int -> Maybe a  -- renvoi l'élément de la liste correspondant à l'index demandé.
getElemList list index =
    if index < 0 || index >= List.length list then
        Nothing
    else
        List.head (List.drop index list)

toString : Http.Error -> String     -- en fonction des erreurs renvois une chaîne de caractère différente
toString erreur = 
  case erreur of 
    Http.BadUrl err -> "Mauvaise URL " ++ err
    Http.Timeout -> "Timeout"
    Http.NetworkError -> "Erreur réseau"
    Http.BadStatus err -> "Mauvais statut" ++ String.fromInt err
    Http.BadBody err -> "BadBody" ++ err 



        
textDatas : (List Datas) -> List (Html Msg)     -- parcours la liste de données et renvoie une liste avec les meanings de chaque data
textDatas datas = 
  case datas of 
    [] -> []
    (x :: xs) -> [li [] ([text "Definitions"] ++ [ul [] (textMeaning x.meanings)])] ++ (textDatas xs)
    

    
textDef : List Definition -> List (Html Msg)          -- parcours la liste de définitions et renvoie une liste avec l'analyse des définitions.
textDef def = 
  case def of
    [] -> []
    (x :: xs) -> [li [] [text x.definition]] ++ (textDef xs)  

textMeaning : List Meaning -> List (Html Msg)           -- parcours la liste des meanings et renvoie une liste avec partOfSpeech et la definition associée
textMeaning meanings = 
  case meanings of
    [] -> []
    (x :: xs) -> [li [] [text x.partOfSpeech]] ++ [ol [] (textDef x.definitions)] ++ (textMeaning xs) 
    
overlay : Model -> List (Html Msg) -> Html Msg
overlay model txt =                                            -- gérer l'affichage sur la page web
  div [style "margin-top" "80px"] 
      [
        if String.toLower model.content == String.toLower model.word then
          div [style "margin-left" "80px", style "font-family" "Copperplate",style "margin-right" "50px"]
           [h1[][text (model.content)]]
        else
          div [style "margin-left" "80px", style "font-family" "Copperplate",style "margin-right" "50px"]
           [h1[][text "Guess the word !!"]],
       div [style "margin-left" "80px", style "font-family" "Copperplate",style "margin-right" "50px"]
         txt
       , div [style "margin-left" "80px",style "font-family" "Helvetica"]
         [ div [style "margin-left" "80px"]
            [input  [style"border" "none",style"border-bottom" "2px solid red",placeholder "écris ta réponse", Html.Attributes.value model.content, onInput Change] [] 
            ,
            if String.toLower model.content == String.toLower model.word then
               div[style "color" "Blue"][h4[] [text "Bravo"]]
            else
               div [style "color" "Blue"][h4[][text ("Tu as entré  " ++ model.content) ]]

                ,button [ style "background-color" "#237689", style "color" "white", style "padding" "15 px 32px", style "text-align" "center" , style "margin" "4px 2px",style "border" "none", style "display" "inline-block" ,onClick Next ] [ text "Next Question" ]
            ]
         ]
      ]
    
  


-- MAIN


main =
  Browser.element{ init = init, update = update, subscriptions = subscriptions, view = viewWord}

