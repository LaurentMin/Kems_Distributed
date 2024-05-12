const cardBackImgPath = 'cards/BACK.png'

const testCardsDraw = [
    { id: 1, value: 2, suit: 'Cloves' },
    { id: 2, value: 3, suit: 'Diamonds' },
    { id: 3, value: 4, suit: 'Hearts' },
    { id: 4, value: 5, suit: 'Spades' }
]

const testCardsHand = [
    { id: 5, value: 6, suit: 'Cloves' },
    { id: 6, value: 9, suit: 'Diamonds' },
    { id: 7, value: 7, suit: 'Hearts' },
    { id: 8, value: 8, suit: 'Spades' }
]

const createAllCardsObject = () => {
    const suits = ['Clubs', 'Diamonds', 'Hearts', 'Spades']
    const cards = []
    let htmlId = 0
    suits.forEach((suit, id) => {
        for (let i = 1; i <= 13; i++) {
            cards.push({ id: htmlId, value: i, suit: suit })
            htmlId++
        }
    })
    return cards
}

const allCards = createAllCardsObject()

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


const playGameButtonElem = document.getElementById('playGame')
const currentGameStatusElem = document.querySelector('.current-status')
const scoreContainerElem = document.querySelector('.header-score-container')
const scoreElem = document.querySelector('.score')
const roundContainerElem = document.querySelector('.header-round-container')
const roundElem = document.querySelector('.round')

const winColor = "green"
const loseColor = "red"
const primaryColor = "black"

let roundNum = 1
let maxRounds = 4
let score = 0

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


function endRound() {
    // TODO:
}

function updateScore() {
    // calculateScore()
    updateStatusElement(scoreElem, "block", primaryColor, `Score <span class='badge'>${score}</span>`)

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


function loadGame() {
    createCards()

    cards = document.querySelectorAll('.card')
    console.log(cards)

    playGameButtonElem.textContent = 'Swap Cards';

    updateStatusElement(scoreContainerElem, "0")
    updateStatusElement(roundContainerElem, "0")

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

// Test
function createCards() {
    testCardsDraw.forEach((cardItem, id) => {
        createCard(cardItem, 'draw', drawPileClassName[id])
    })

    testCardsHand.forEach((cardItem, id) => {
        createCard(cardItem, 'hand', handPileClassName[id])
    })

    createCard({ id: 9, value: 10, suit: 'Spades' }, 'deck', deckClassName)
}


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
            newTurn()
            break;

    }

}

// Back function
function newTurn() {
    console.log('New Turn')
}

function swapCards(drawCard, handCard) {
    console.log('Swaxpping Cards')
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
function getObjectFromJSON(json) {
    return JSON.parse(json)
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

loadGame()