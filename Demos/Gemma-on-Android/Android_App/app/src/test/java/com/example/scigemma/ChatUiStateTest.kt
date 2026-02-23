package com.example.scigemma

import org.junit.Test
import java.util.UUID

class ChatUiStateTest {

    @Test
    fun benchmarkAppendMessage() {
        val messageCount = 5000
        val messages = (1..messageCount).map {
            ChatMessage(
                id = UUID.randomUUID().toString(),
                message = "Message $it",
                author = if (it % 2 == 0) USER_PREFIX else MODEL_PREFIX
            )
        }

        val uiState = GemmaUiState(messages)

        // Target messages: first, middle, last
        val firstId = messages.first().id
        val middleId = messages[messageCount / 2].id
        val lastId = messages.last().id

        val iterations = 1000

        val start = System.nanoTime()
        for (i in 0 until iterations) {
            uiState.appendMessage(firstId, " a", false)
            uiState.appendMessage(middleId, " b", false)
            uiState.appendMessage(lastId, " c", false)
        }
        val end = System.nanoTime()

        println("Time taken for $iterations iterations with $messageCount messages: ${(end - start) / 1_000_000.0} ms")
    }

    @Test
    fun testGemmaUiStateAppendMessageCorrectness() {
        val id1 = "id1"
        val id2 = "id2"
        val messages = listOf(
            ChatMessage(id = id1, message = "Hello", author = USER_PREFIX),
            ChatMessage(id = id2, message = "Hi", author = MODEL_PREFIX)
        )
        val uiState = GemmaUiState(messages)

        // Append to id2
        uiState.appendMessage(id2, " there", false)

        val updatedMessages = uiState.messages
        // messages are reversed
        val updatedMessage2 = updatedMessages.find { it.id == id2 }!!

        // "Hi there"
        assert(updatedMessage2.message == "Hi there")
    }

    @Test
    fun testChatUiStateAppendMessageCorrectness() {
        val id1 = "id1"
        val id2 = "id2"
        val messages = listOf(
            ChatMessage(id = id1, message = "Hello", author = USER_PREFIX),
            ChatMessage(id = id2, message = "Hi", author = MODEL_PREFIX)
        )
        val uiState = ChatUiState(messages)

        uiState.appendMessage(id2, " there", false)

        val updatedMessages = uiState.messages
        val updatedMessage2 = updatedMessages.find { it.id == id2 }!!
        assert(updatedMessage2.message == "Hi there")
    }

    @Test
    fun testAddMessageCorrectness() {
         val uiState = GemmaUiState(emptyList())
         val id = uiState.addMessage("New message", USER_PREFIX)

         // Verify it exists
         val messages = uiState.messages
         assert(messages.any { it.id == id })

         // Verify map logic by appending to it
         uiState.appendMessage(id, " appended", false)

         val updated = uiState.messages.find { it.id == id }!!
         // "New message appended"
         assert(updated.message.contains("New message appended"))
    }

    @Test
    fun testCreateLoadingMessageCorrectness() {
         val uiState = GemmaUiState(emptyList())
         val id = uiState.createLoadingMessage()

         val messages = uiState.messages
         assert(messages.any { it.id == id && it.isLoading })

         // Append to it
         uiState.appendMessage(id, "Loaded", true)

         val updated = uiState.messages.find { it.id == id }!!
         assert(!updated.isLoading)
         assert(updated.message == "Loaded")
    }
}
