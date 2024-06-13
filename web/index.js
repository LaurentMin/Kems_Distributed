// Get the modal
var modal = document.getElementById("myModal");

// Get the button that opens the modal
var btn = document.getElementById("openModal");

// Get the <span> element that closes the modal
var span = document.getElementsByClassName("close")[0];

// When the user clicks the button, open the modal 
btn.onclick = function () {
    modal.style.display = "block";
}

var hardResetButton = document.getElementById('reset-game');
hardResetButton.onclick = function () {
    actionObject = {
        action: 'ResetGame',
    }
    sendAction(actionObject)
    modal.style.display = "none";
}

// When the user clicks on <span> (x), close the modal
span.onclick = function () {
    modal.style.display = "none";
}

var scoreModalContent = document.getElementById('scoreModal');

var span2 = document.getElementsByClassName("close")[1];
span2.onclick = function () {
    scoreModalContent.style.display = "none";
}
window.onclick = function (event) {
    if (event.target == scoreModalContent) {
        scoreModalContent.style.display = "none";
    }
    if (event.target == modal) {
        modal.style.display = "none";
    }
}

// When the user clicks anywhere outside of the modal, close it
window.onclick = function (event) {
    if (event.target == modal) {
        modal.style.display = "none";
    }
}



// Function to update the score modal
function updateScoreModal(scores) {
    const scoreElem = document.getElementById('scoreModalContent');

    // Clear the current content
    scoreElem.innerHTML = '';

    // Add the scores
    maxScore = 0;
    winnerId = 0;
    gameState.scores.forEach((score, id) => {
        const scoreElement = document.createElement('p');
        scoreElement.textContent = `Player ${id + 1}: ${score}`;
        if (score > maxScore) {
            maxScore = score;
            winnerId = id + 1;
        }
        scoreElem.appendChild(scoreElement);
    });
    // Display current winner
    const winnerElement = document.getElementById('winner');
    winnerElement.textContent = `Player ${winnerId} is winning!`;

    scoreModalContent.style.display = "block";

}

const drawPileClassName = [
    '.card-pos-a',
    '.card-pos-b',
    '.card-pos-c',
    '.card-pos-d'
]

const handPileClassName = [
    '.card-pos-a-player',
    '.card-pos-b-player',
    '.card-pos-c-player',
    '.card-pos-d-player'
]

const deckClassName = ".card-deck"
const discardClassName = ".card-discard"




const connectedElem = document.getElementById('connected')
const playGameButtonElem = document.getElementById('playGame')
const playGameButtonContainerElem = document.querySelector('.header-button-container')
const currentGameStatusElem = document.querySelector('.current-status')
const scoreContainerElem = document.querySelector('.header-score-container')
const scoreElem = document.querySelector('.score')
const roundContainerElem = document.querySelector('.header-round-container')
const roundElem = document.querySelector('.round')

const winColor = "green"
const loseColor = "red"
const primaryColor = "black"



const createAllCardsObject = () => {
    const suits = ['Clubs', 'Diamonds', 'Hearts', 'Spades']
    const cards = []
    let htmlId = 0
    suits.forEach((suit, id) => {
        cards.push({ id: htmlId, value: 'A', suit: suit })
        htmlId++
        for (let i = 2; i <= 10; i++) {
            cards.push({ id: htmlId, value: `${i}`, suit: suit })
            htmlId++
        }
        cards.push({ id: htmlId, value: 'J', suit: suit })
        htmlId++
        cards.push({ id: htmlId, value: 'Q', suit: suit })
        htmlId++
        cards.push({ id: htmlId, value: 'K', suit: suit })
        htmlId++
    })
    return cards
}

const allCards = createAllCardsObject()

// Web socket connection
var ws;
let selectedPlayer = '1'; // Default to player 1

// Listen for changes on the radio buttons
document.querySelectorAll('input[name="player"]').forEach((elem) => {
    elem.addEventListener('change', function () {
        selectedPlayer = this.value;
    });
});

document.getElementById("connect").onclick = function (evt) {
    if (ws) {
        return false;
    }

    var host = "localhost";
    let port
    switch (selectedPlayer) {
        case '1':
            port = 4444;
            break;
        case '2':
            port = 5555;
            break;
        case '3':
            port = 5000;
            break;
    }

    try {
        ws = new WebSocket("ws://" + host + ":" + port + "/ws");
    }
    catch (err) {
        console.error("Error connecting to websocket", err)
        return false;
    }

    ws.onopen = function (evt) {
        console.log("Connection open ...");
        connectedElem.textContent = "Connected";
        ws.send("Hello Server!!")
    };

    ws.onclose = function (evt) {
        connectedElem.textContent = "";
        ws = null;
    }

    ws.onmessage = function (evt) {
        console.log("Received message: " + evt.data);
        gameState = JSON.parse(evt.data)
        updateGame(gameState)
    }

    ws.onerror = function (evt) {
        console.log("Connection error ...");
    }

    return false;
}

document.getElementById("close").onclick = function (evt) {
    if (!ws) {
        return false;
    }
    ws.close();
    return false;
}

const cardBackImgPath = 'cards/BACK.png'

const getCardWithHTMLId = (cardItem) => {
    test = allCards.find(card => card.value === cardItem.value && card.suit === cardItem.suit)
    if (!test) {
        console.error("ERROR Card not found", cardItem)
    }
    return allCards.find(card => card.value === cardItem.value && card.suit === cardItem.suit)
}

const updateGame = (gameState) => {
    console.log("Updating game state", gameState)
    gameState.drawPile.forEach((cardItem, id) => {
        cardItemWithHTMLId = getCardWithHTMLId(cardItem)
        createCard(cardItemWithHTMLId, 'draw', drawPileClassName[id])
    })

    gameState.hand.forEach((cardItem, id) => {
        cardItemWithHTMLId = getCardWithHTMLId(cardItem)
        createCard(cardItemWithHTMLId, 'hand', handPileClassName[id])
    })

    //Create False card for deck
    createCard({ id: 99, value: 1, suit: 'Spades' }, 'deck', deckClassName)

    // discardPile might not be in the gameState
    if (gameState.discardPile) {
        cardItem = gameState.discardPile[0]
        cardItemWithHTMLId = getCardWithHTMLId(cardItem)
        createCard(cardItemWithHTMLId, 'discard', discardClassName)
    }
    if (gameState.newRound) {
        updateScoreModal(gameState.scores)
    }

    if (gameState.potentialWinner) {
        playGameButtonElem.style.display = "inline-block";
        if (gameState.potentialWinner === gameState.playerId) {
            playGameButtonElem.textContent = 'Kems';
            playGameButtonElem.addEventListener('click', kems);
        } else {
            playGameButtonElem.textContent = 'Counter Kems';
            playGameButtonElem.addEventListener('click', counterKems);
        }
    } else {
        playGameButtonElem.textContent = 'Swap Cards';
    }

    playerScore = gameState.scores[gameState.playerId - 1]
    console.log("Player score", playerScore, gameState.round)

    updateStatusElement(scoreElem, "block", primaryColor, `Score&nbsp;<span class='badge'>${playerScore}</span>`)
    updateStatusElement(roundElem, "block", primaryColor, `Round&nbsp;<span class='badge'>${gameState.round}</span>`)
}

function sendAction(action) {
    if (!ws) {
        console.error('Websocket not connected')
        return false;
    }
    console.log('Sending action', action)
    ws.send(JSON.stringify(actionObject))
}

// game Action
function newTurn() {
    actionObject = {
        action: 'NextTurn',
    }
    sendAction(actionObject)
}

function swapCards(drawCardHtml, handCardHtml) {
    console.log('Swapping Cards')
    drawCard = allCards.find(card => card.id === parseInt(drawCardHtml.id))
    handCard = allCards.find(card => card.id === parseInt(handCardHtml.id))
    console.log('Draw Card', drawCard, drawCardHtml)
    console.log('Hand Card', handCard, handCardHtml)
    actionObject = {
        action: 'SwapCards',
        drawCardValue: `${drawCard.value}`,
        drawCardSuit: drawCard.suit,
        handCardValue: `${handCard.value}`,
        handCardSuit: handCard.suit,
    }
    sendAction(actionObject)
}

function kems() {
    console.log('Kems')
    actionObject = {
        action: 'Kems',
    }
    sendAction(actionObject)
    playGameButtonElem.style.display = "none"
}

function counterKems() {
    console.log('Counter Kems')
    actionObject = {
        action: 'ContreKems',
    }
    sendAction(actionObject)
    playGameButtonElem.style.display = "none"
}


const getCardValueForFileName = (card) => {
    if (card.value === 0) {
        return ''
    }
    else if (card.value === 1) {
        return 'A'
    } else if (card.value === 11) {
        return 'J'
    }
    else if (card.value === 12) {
        return 'Q'
    }
    else if (card.value === 13) {
        return 'K'
    }
    else {
        return card.value
    }
}

let cardPositions = []


function endRound() {
    // TODO:
}

function updateStatusElement(elem, display, color, innerHTML) {
    elem.style.display = display

    if (arguments.length > 2) {
        elem.style.color = color
        elem.innerHTML = innerHTML
    }

}

function outputChoiceFeedBack(hit) {
    if (hit) {
        updateStatusElement(currentGameStatusElem, "block", winColor, "Hit!! - Well Done!! :)")
    }
    else {
        updateStatusElement(currentGameStatusElem, "block", loseColor, "Missed!! :(")
    }
}


// Swap cards logic
let drawPileSelectedCard = null;
let handPileSelectedCard = null;


function swapButtonHandler() {
    if (drawPileSelectedCard && handPileSelectedCard) {
        swapCards(drawPileSelectedCard, handPileSelectedCard);

        drawPileSelectedCard.classList.remove('selected');
        handPileSelectedCard.classList.remove('selected');
        drawPileSelectedCard = null;
        handPileSelectedCard = null;
        playGameButtonElem.style.display = "none"
    }
}




function initializeNewRound() {
    roundNum++

    updateStatusElement(currentGameStatusElem, "block", primaryColor, "Shuffling...")

    updateStatusElement(roundElem, "block", primaryColor, `Round <span class='badge'>${roundNum}</span>`)

}


function flipCard(card, flipToBack) {
    const innerCardElem = card.firstChild

    if (flipToBack) {
        innerCardElem.classList.add('flip-it')
    }
    else if (innerCardElem.classList.contains('flip-it')) {
        innerCardElem.classList.remove('flip-it')
    }

}


/* <div class="card">
<div class="card-inner">
    <div class="card-front">
        <img src="/images/card-JackClubs.png" alt="" class="card-img">
    </div>
    <div class="card-back">
        <img src="/images/card-back-Blue.png" alt="" class="card-img">
    </div>
</div>
</div> */

function createCard(cardItem, pile, cardPositionClassName) {

    //create div elements that make up a card
    const cardElem = createElement('div')
    const cardInnerElem = createElement('div')
    const cardFrontElem = createElement('div')
    const cardBackElem = createElement('div')

    //create front and back image elements for a card
    const cardFrontImg = createElement('img')
    const cardBackImg = createElement('img')

    //add class and id to card element
    addClassToElement(cardElem, 'card')
    addIdToElement(cardElem, cardItem.id)

    //add class to inner card element
    addClassToElement(cardInnerElem, 'card-inner')

    //add class to front card element
    addClassToElement(cardFrontElem, 'card-front')

    //add class to back card element
    addClassToElement(cardBackElem, 'card-back')

    //add src attribute and appropriate value to img element - back of card
    addSrcToImageElem(cardBackImg, cardBackImgPath)

    //add src attribute and appropriate value to img element - front of card
    value = getCardValueForFileName(cardItem)
    imgPath = `/cards/${value}-${cardItem.suit[0]}.png`
    addSrcToImageElem(cardFrontImg, `cards/${value}-${cardItem.suit[0]}.png`)

    //assign class to back image element of back of card
    addClassToElement(cardBackImg, 'card-img')

    //assign class to front image element of front of card
    addClassToElement(cardFrontImg, 'card-img')

    //add front image element as child element to front card element
    addChildElement(cardFrontElem, cardFrontImg)

    //add back image element as child element to back card element
    addChildElement(cardBackElem, cardBackImg)

    //add front card element as child element to inner card element
    addChildElement(cardInnerElem, cardFrontElem)

    //add back card element as child element to inner card element
    addChildElement(cardInnerElem, cardBackElem)

    //add inner card element as child element to card element
    addChildElement(cardElem, cardInnerElem)

    //add card element as child element to appropriate grid cell
    const cardPosElem = document.querySelector(cardPositionClassName)
    addChildElement(cardPosElem, cardElem)


    initializeCardPositions(cardElem)

    attatchClickEventHandlerToCard(cardElem, pile)

    if (pile === 'deck') {
        flipCard(cardElem, true)
    }
}


function attatchClickEventHandlerToCard(card, pile) {
    switch (pile) {
        case 'draw':
            card.addEventListener('click', () => {
                console.log('Hand pile card clicked', card);
                if (drawPileSelectedCard) {
                    console.log('Previous card deselected',);
                    drawPileSelectedCard.classList.remove('selected');
                    if (drawPileSelectedCard === card) {
                        drawPileSelectedCard = null;
                        playGameButtonElem.style.display = "none"
                        return;
                    }
                }
                card.classList.add('selected');
                drawPileSelectedCard = card;
                console.log('Draw pile card selected', drawPileSelectedCard);

                if (drawPileSelectedCard && handPileSelectedCard) {
                    // Remove old event listener
                    playGameButtonElem.removeEventListener('click', swapButtonHandler);

                    playGameButtonElem.textContent = 'Swap Cards';
                    playGameButtonElem.style.display = "inline-block";
                    // Add new event listener
                    playGameButtonElem.addEventListener('click', swapButtonHandler);
                } else {
                    playGameButtonElem.style.display = "none"
                }
            });
            break;
        case 'hand':
            card.addEventListener('click', () => {
                console.log('Hand pile card clicked', card);
                if (handPileSelectedCard) {
                    console.log('Previous card deselected',);
                    handPileSelectedCard.classList.remove('selected');
                    if (handPileSelectedCard === card) {
                        handPileSelectedCard = null;
                        playGameButtonElem.style.display = "none"
                        return;
                    }
                }
                card.classList.add('selected');
                handPileSelectedCard = card;
                console.log('Hand pile card selected', handPileSelectedCard);

                if (drawPileSelectedCard && handPileSelectedCard) {
                    // Remove old event listener
                    playGameButtonElem.removeEventListener('click', swapButtonHandler);

                    playGameButtonElem.textContent = 'Swap Cards';
                    playGameButtonElem.style.display = "inline-block";
                    // Add new event listener
                    playGameButtonElem.addEventListener('click', swapButtonHandler);
                } else {
                    playGameButtonElem.style.display = "none"
                }
            });
            break;
        case 'deck':
            card.addEventListener('click', () => {
                newTurn()
            });
            break;

    }

}

function createElement(elemType) {
    return document.createElement(elemType)
}
function initializeCardPositions(card) {
    cardPositions.push(card.id)
}
function addClassToElement(elem, className) {
    elem.classList.add(className)
}
function addIdToElement(elem, id) {
    elem.id = id
}

function addSrcToImageElem(imgElem, src) {
    imgElem.src = src
}

function addChildElement(parentElem, childElem) {
    parentElem.appendChild(childElem)
}


//local storage functions
function getSerializedObjectAsJSON(obj) {
    return JSON.stringify(obj)
}
function updateLocalStorageItem(key, value) {
    localStorage.setItem(key, value)
}
function removeLocalStorageItem(key) {
    localStorage.removeItem(key)
}
function getLocalStorageItemValue(key) {
    return localStorage.getItem(key)
}
