/*
This object will track mouse movement events from an given html5 element, and relay them
with POST messages to a an http server. (asynchronously).

The constructor takes care of registering itself as an event handler for the html element you
specify.

It takes some steps to optimise the communications assuming that the incoming mouse event stream
would far exceed how fast we can send them to the server. Firstly it ignores small mouse movements.
Secondly it throttles its sending of messages to match the rate at which the server can receive
and reply to them. It does this using a "lossy" back-pressure strategy, whereby a maximum of <N>
POSTS are in flight at any one time. Mouse movements received while the pipe has backed up are
"spilled", but we keep track of the most recently received only, and arranging for that one to be
flushed as soon as the POST's reply callback detects room in the pipe. In the steady
throttled  state, this means we are sending events at the maximum rate our communication with the
server can handle, but we get the benefit of sending messages concurrently without waiting for
the previous message's round trip.
*/

// Factory method for the object.
function mouseEventTransmitter(mouseEventSource) {
    t = {};
    // Constants
    t.rect = mouseEventSource.getBoundingClientRect();
    t.left = t.rect.left;
    t.top = t.rect.top;
    t.width = t.rect.right - t.rect.left;
    t.height = t.rect.bottom - t.rect.top;
    t.serverUrl = "/mouse";

    t.negligibleMovementSize = 5;
    t.maxInFlight = 3; // The <N> in the description above.

    // State
    t.prevX = 1e6; // An out of band sentinel value.
    t.prevY = 1e6;
    t.numberOfMessagesInFlight = 0;
    t.mostRecentUnSentPosition = null;

    // Functions

    // Handler to receive mouse move events from the browser
    t.handleMoveEvent= function(event) {
        mouseX = event.clientX - t.left;
        mouseY = event.clientY - t.top;
        // Optimise out, and ignore moves that are too small (or zero) while still in integer
        // coordinates
        if (t.movementTooSmall(mouseX, mouseY)) {
            return;
        }
        t.prevX = mouseX; // used to test for too small
        t.prevY = mouseY;

        normalisedX = mouseX / t.width;
        normalisedY = mouseY / t.height;

        // Send the position to the server, in accordance with lossy back pressure
        // strategy.

        // If the pipe is currently full (the general case), all we do is note the most
        // recently received position, to flush later.
        if (t.numberOfMessagesInFlight == t.maxInFlight) {
            t.mostRecentUnSentPosition = {x: normalisedX, y: normalisedY};
        }
        else {
            // Seems we are at liberty to forward the movement, so we do so.
            t.mostRecentUnSentPosition = null;
            t.postPositionToServerAsync(normalisedX, normalisedY);
        }
    };

    // Utility function to assess if mouse movement is very small or zero
    t.movementTooSmall = function(x, y) {
          if (Math.abs(x - t.prevX) > t.negligibleMovementSize)
              return false;
          if (Math.abs(y - t.prevY) > t.negligibleMovementSize)
              return false;
          return true;
     };

    // Utility to post an x / y position to the server.
    t.postPositionToServerAsync = function(x, y) {
        payload = {X: normalisedX, Y: normalisedY};
        callbackFn = t.handlerForServerReply;
        t.numberOfMessagesInFlight += 1;
        $.post(t.serverUrl, payload, callbackFn);
    }

    // Handler to accept replies from the server.
    t.handlerForServerReply = function(reply) {
        // Receiving a reply triggers our semaphore maintenance. We are not interested
        // in the content of the reply.
        t.numberOfMessagesInFlight -= 1;
        // There must be some room in the pipe at this point, having just decremented
        // the counter, so we seize the opportunity to flush the most recently received
        // position that has not yet been sent, if there is one.
        if (t.mostRecentUnSentPosition != null) {
            x = t.mostRecentUnSentPosition.x;
            y = t.mostRecentUnSentPosition.y;
            t.mostRecentUnSentPosition = null; // spent
            t.postPositionToServerAsync(x, y);
        }
    }

    // Ask the mouse move event emitter to send the events to a method in this object.
    mouseEventSource.onmousemove = t.handleMoveEvent;
    return t;
}