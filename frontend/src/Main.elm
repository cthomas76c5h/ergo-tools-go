module Main exposing (main)

import Browser
import Html exposing (..)
import Html.Attributes as A
import Html.Events as E
import Http
import Json.Decode as D
import Json.Encode as JE


-- MODEL

type alias MsgItem =
    { role : String
    , text : String
    }

type alias Model =
    { server : String
    , tenantId : String
    , sessionId : String
    , message : String
    , history : List MsgItem
    , lead : D.Value
    }


init : () -> ( Model, Cmd Msg )
init _ =
    ( { server = "http://localhost:8000"
      , tenantId = "dev"
      , sessionId = genSession 0
      , message = ""
      , history = []
      , lead = JE.null
      }
    , Cmd.none
    )


-- MESSAGES

type Msg
    = SetServer String
    | SetTenant String
    | SetSession String
    | NewSession
    | SetMessage String
    | Send
    | GotHttp (Result Http.Error ChatReply)


-- UPDATE

update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        SetServer s ->
            ( { model | server = s }, Cmd.none )

        SetTenant t ->
            ( { model | tenantId = t }, Cmd.none )

        SetSession s ->
            ( { model | sessionId = s }, Cmd.none )

        NewSession ->
            ( { model | sessionId = genSession (List.length model.history)
              , history = []
              , lead = JE.null
              }
            , Cmd.none
            )

        SetMessage m ->
            ( { model | message = m }, Cmd.none )

        Send ->
            let
                trimmed =
                    String.trim model.message

                model1 =
                    if trimmed == "" then
                        model
                    else
                        { model
                            | message = ""
                            , history = model.history ++ [ MsgItem "user" trimmed ]
                        }
            in
            if trimmed == "" then
                ( model, Cmd.none )

            else
                ( model1
                , httpSend model1.server model1.tenantId model1.sessionId trimmed
                )

        GotHttp result ->
            case result of
                Ok rep ->
                    ( { model
                        | history = model.history ++ [ MsgItem "assistant" rep.reply ]
                        , lead = rep.lead
                      }
                    , Cmd.none
                    )

                Err e ->
                    ( { model
                        | history = model.history ++ [ MsgItem "assistant" ("⚠️ " ++ httpErr e) ]
                      }
                    , Cmd.none
                    )


-- VIEW

view : Model -> Html Msg
view model =
    div [ A.class "container" ]
        [ h2 [] [ text "Ergo Tools Chat Tester (Elm) — HTTP" ]
        , div [ A.class "row" ]
            [ labeledInput "Server" model.server SetServer
            , labeledInput "Tenant" model.tenantId SetTenant
            , div []
                [ labeledInput "Session ID" model.sessionId SetSession
                , button [ A.class "btn", E.onClick NewSession, A.style "margin-left" "8px" ] [ text "New" ]
                ]
            ]
        , div [ A.class "card", A.style "margin-top" "12px" ]
            [ h3 [] [ text "Chat" ]
            , div [ A.style "height" "50vh", A.style "overflow" "auto", A.class "card", A.style "background" "#fff" ]
                [ if List.isEmpty model.history then
                    div [ A.class "hint" ] [ text "Try: ‘Car accident yesterday. My email is john dot doe at gmail dot com.’" ]
                  else
                    div [] (List.map viewMsg model.history)
                ]
            , div [ A.style "display" "flex", A.style "gap" "8px", A.style "margin-top" "8px" ]
                [ textarea
                    [ A.class "input"
                    , A.style "flex" "1"
                    , A.rows 2
                    , A.placeholder "Type a message and press Enter"
                    , A.value model.message
                    , E.onInput SetMessage
                    , onEnter Send
                    ]
                    []
                , button [ A.class "btn", E.onClick Send ] [ text "Send" ]
                ]
            ]
        , div [ A.class "card", A.style "margin-top" "12px" ]
            [ h3 [] [ text "Lead Inspector" ]
            , pre [] [ text (JE.encode 2 model.lead) ]
            ]
        , styles
        ]


viewMsg : MsgItem -> Html msg
viewMsg m =
    let
        cls =
            if m.role == "user" then
                "bubble user"
            else
                "bubble bot"
    in
    div [ A.class cls ] [ text m.text ]


labeledInput : String -> String -> (String -> msg) -> Html msg
labeledInput label value toMsg =
    div []
        [ span [ A.class "hint" ] [ text label ]
        , input [ A.class "input", A.value value, E.onInput toMsg ] []
        ]


styles : Html msg
styles =
    node "style" []
        [ text """
        body{font-family:ui-sans-serif, system-ui, -apple-system, Segoe UI, Roboto, Inter, sans-serif;background:#f8fafc;color:#0f172a;margin:0}
        .container{max-width:1100px;margin:24px auto;padding:0 16px}
        .card{background:#f1f5f9;border:1px solid #e2e8f0;border-radius:16px;padding:16px;box-shadow:0 1px 2px rgba(0,0,0,0.05)}
        .btn{background:#0284c7;color:#fff;border:none;border-radius:12px;padding:10px 14px;font-weight:600;cursor:pointer}
        .btn:disabled{opacity:.6;cursor:default}
        .input, textarea{border:1px solid #cbd5e1;border-radius:12px;padding:10px;background:#fff}
        .hint{font-size:12px;color:#475569}
        .bubble{display:inline-block;max-width:85%;border-radius:16px;padding:8px 12px;margin:4px 0}
        .user{background:#0369a1;color:#fff;float:right;clear:both}
        .bot{background:#e2e8f0;color:#0f172a;float:left;clear:both}
        """ ]


-- HTTP

type alias ChatReply =
    { reply : String
    , lead : D.Value
    }

chatReplyDecoder : D.Decoder ChatReply
chatReplyDecoder =
    D.map2 ChatReply
        (D.field "reply" D.string)
        (D.field "lead" D.value)


httpSend : String -> String -> String -> String -> Cmd Msg
httpSend server tenant session text =
    let
        body =
            JE.object
                [ ( "tenant_id", JE.string tenant )
                , ( "session_id", JE.string session )
                , ( "message", JE.string text )
                ]
    in
    Http.post
        { url = server ++ "/chat"
        , body = Http.jsonBody body
        , expect = Http.expectJson GotHttp chatReplyDecoder
        }


httpErr : Http.Error -> String
httpErr e =
    case e of
        Http.BadUrl u ->
            "Bad URL: " ++ u

        Http.Timeout ->
            "Request timed out"

        Http.NetworkError ->
            "Network error"

        Http.BadStatus _ ->
            "Bad status"

        Http.BadBody m ->
            "Bad body: " ++ m


-- HELPERS

onEnter : msg -> Html.Attribute msg
onEnter msg =
    let
        decode =
            D.field "key" D.string
                |> D.andThen (\k -> if k == "Enter" then D.succeed msg else D.fail "")
    in
    E.on "keydown" decode


genSession : Int -> String
genSession n =
    let
        base = "xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx"

        step i ch =
            if ch == 'x' then
                hex ((i + n) % 16)
            else if ch == 'y' then
                hex (((i + n) % 4) + 8)
            else
                String.fromChar ch
    in
    String.indexedMap step base


hex : Int -> String
hex i =
    let
        chars =
            "0123456789abcdef"

        get k =
            case String.uncons (String.dropLeft k chars) of
                Just (c, _) ->
                    String.fromChar c

                Nothing ->
                    "0"
    in
    get i


-- PROGRAM

main : Program () Model Msg
main =
    Browser.element
        { init = init
        , update = update
        , view = view
        , subscriptions = \_ -> Sub.none
        }
