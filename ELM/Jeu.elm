

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

type State = Failure String | Loading | Success String

type alias Datas = { word : String, meanings : List Meaning}
type alias Meaning = {partOfSpeech : String, definitions : List Definition}    
type alias Definition = {definition : String}

type alias Model = 
    { http : State
    , content : String 
    , lWords : List String
    , word : String
    , jSon : State
    , datas : List Datas
    }

init : () -> (Model, Cmd Msg)             
init _ =
  ( Model Loading "" [] "" Loading []
  , Http.get
      { url = " http://localhost:8000/monFichier.txt"
      , expect = Http.expectString Words
      }
  )

-- UPDATE


type Msg
  = Words (Result Http.Error String)
  | RandomWord Int
  | GotJson (Result Http.Error (List Datas))
  | Change String

  
update : Msg -> Model -> (Model, Cmd Msg)
update msg model =
  case msg of
    Words result ->
      case result of
        Ok words ->
          ({ model | lWords = String.words words , http = Success "" } , randomWord model)

        Err error ->
          ({model | http = Failure (toString error)}, Cmd.none)
          
    RandomWord index -> case (getElemList model.lWords index) of
                                Nothing -> (model, randomWord model)
                                Just x -> ({ model | word = x }, Http.get {url = ("https://api.dictionaryapi.dev/api/v2/entries/en/" ++ x)  , expect = Http.expectJson GotJson lDatasDecoder})
     
    GotJson result -> case result of
                            Ok data-> ({ model | jSon = Success "" , datas = data} , Cmd.none)

                            Err error -> ({ model | jSon = Failure (toString error) } , randomWord model)   

    Change newContent -> ({model | content = newContent}, Cmd.none)
    
randomWord : Model -> Cmd Msg
randomWord model = Random.generate RandomWord (Random.int 1 (List.length model.lWords)) -- erreur potentielle ici 


-- VIEW
viewWord : Model -> Html Msg
viewWord model =
  case model.http of
    Loading -> text "Chargement de la page"
    
    Failure error -> text ("Chargement impossible : " ++ error) 

    Success good -> overlay model (
      case model.jSon of
        Success veryGood -> [div [] (List.map viewData model.datas)]
        Loading -> [text "Chargement"]
        Failure error -> [text ("échec " ++ error)] )

viewData : Datas -> Html Msg
viewData data =
    div []
        [ h3 [] [text data.word]
        , ul [] (List.map viewMeaning data.meanings)
        ]

viewMeaning : Meaning -> Html Msg
viewMeaning meaning =
    div []
        [ h4 [] [text meaning.partOfSpeech]
        , ul [] (List.map viewDefinition meaning.definitions)
        ]

viewDefinition : Definition -> Html Msg
viewDefinition definition =
    li [] [text definition.definition]





-- SUBSCRIPTIONS


subscriptions : Model -> Sub Msg
subscriptions model =
  Sub.none




          
-- JSON 

lDatasDecoder : Decoder (List Datas)
lDatasDecoder = Json.Decode.list datasDecoder
definitionDecoder : Decoder Definition
definitionDecoder = Json.Decode.map Definition (field "definition" string)
datasDecoder : Decoder Datas
datasDecoder = map2 Datas (field "word" string)(field "meanings" <| Json.Decode.list meaningDecoder)
meaningDecoder : Decoder Meaning
meaningDecoder = map2 Meaning (field "partOfSpeech" string)(field "definitions" <| Json.Decode.list definitionDecoder)


-- HELPERS
getElemList : List a -> Int -> Maybe a
getElemList list index =
    if index < 0 || index >= List.length list then
        Nothing
    else
        List.head (List.drop index list)

toString : Http.Error -> String 
toString erreur = 
  case erreur of 
    Http.BadUrl err -> "BadUrl" ++ err
    Http.Timeout -> "Timeout"
    Http.NetworkError -> "NetworkError"
    Http.BadStatus err -> "BadStatus" ++ String.fromInt err
    Http.BadBody err -> "BadBody" ++ err 



        
textDatas : (List Datas) -> List (Html Msg)
textDatas datas = 
  case datas of 
    [] -> []
    (x :: xs) -> [li [] ([text "Definitions"] ++ [ul [] (textMeaning x.meanings)])] ++ (textDatas xs)
    

    
textDef : List Definition -> List (Html Msg)
textDef def = 
  case def of
    [] -> []
    (x :: xs) -> [li [] [text x.definition]] ++ (textDef xs)  

textMeaning : List Meaning -> List (Html Msg)
textMeaning meanings = 
  case meanings of
    [] -> []
    (x :: xs) -> [li [] [text x.partOfSpeech]] ++ [ol [] (textDef x.definitions)] ++ (textMeaning xs) 
    
overlay : Model -> List (Html Msg) -> Html Msg
overlay model txt = 
  div [] 
      [
       div [style "text-align" "left"]
         txt
       , div [style "text-align" "center"]
         [ div []
            [input [placeholder "écris ta réponse", Html.Attributes.value model.content, onInput Change] [] 
            ,
            if String.toLower model.content == String.toLower model.word then
               div[style "color" "Blue" ] [text "Bravo"]
            else
               div [] [text ("Tu as entré  " ++ model.content) ]
            ]
         ]
      ]
    


-- MAIN


main =
  Browser.element{ init = init, update = update, subscriptions = subscriptions, view = viewWord}
