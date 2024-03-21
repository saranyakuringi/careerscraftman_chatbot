#indel.html 

<!DOCTYPE html>
<html>
<head>
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>CareerCraftMan Bot</title>
    <link rel="stylesheet" href="style.css">
    <script src="https://kit.fontawesome.com/854e65de19.js" crossorigin="anonymous"></script>
</head>
<body>
    <div class="container">
        <div class="form-box">
            <h1 id="title">Sign Up</h1>
            <form>
                <div class="input-group">
                    <div class="input-field" id="firstName">
                        <i class="fa-solid fa-user"></i>
                        <input type="text" placeholder="First Name"> 
                    </div>

                    <div class="input-field" id="lastName">
                        <i class="fa-solid fa-user"></i>
                        <input type="text" placeholder="Last Name"> 
                    </div>

                    <div class="input-field">
                        <i class="fa-solid fa-user"></i>
                        <input type="Username" placeholder="User Name"> 
                    </div>

                    <div class="input-field">
                        <i class="fa-solid fa-lock"></i>
                        <input type="password" placeholder="Password"> 
                    </div>
                    <p>Lost password <a href="#">Click Here!</a></p>
                </div>
                <div class="btn-field">
                    <button type="button" id="signupBtn">Sign up</button>
                    <button type="button" id="signinBtn" class="disable">Sign in</button>
                </div>
            </form>
        </div>
    </div>
<script>
    let signupBtn=document.getElementById("signupBtn");
    let signinBtn=document.getElementById("signinBtn");
    let nameField=document.getElementById("nameField");
    let title=document.getElementById("title");

    signinBtn.onclick=function(){
        firstName.style.maxHeight="0";
        lastName.style.maxHeight="0";
        title.innerHTML="Sign In";
        signupBtn.classList.add("disable");
        signinBtn.classList.remove("disable");
    }

    signupBtn.onclick=function(){
        firstName.style.maxHeight="60px";
        lastName.style.maxHeight="60px";
        title.innerHTML="Sign Up";
        signupBtn.classList.remove("disable");
        signinBtn.classList.add("disable");
    }

</script>
        
</body>
</html>


#signin.html

<!DOCTYPE html>
<html>

<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Sign In</title>
    <link rel="stylesheet" href="/static/style.css">
</head>

<body>
    <div class="container">
        <div class="form-box">
            <h1>Sign In</h1>
            <form id="signinForm" method="POST" action="/signin">
                <div class="input-field">
                    <i class="fa-solid fa-user"></i>
                    <input type="text" id="username" name="username" placeholder="Username" required>
                </div>
                <div class="input-field">
                    <i class="fa-solid fa-lock"></i>
                    <input type="password" id="password" name="password" placeholder="Password" required>
                </div>
                <div id="error-message" class="error-message">{{ .error }}</div>
                <button type="submit">Sign In</button>
            </form>
        </div>
    </div>
</body>

</html>


router.POST("/signin", func(c *gin.Context) {
		fmt.Println(c)
		//fmt.Println(c.Request.Form)

		fmt.Println("Form Data:", c.Request.Form)

		fmt.Println("Post Form Data:", c.Request.PostForm)

		//fmt.Println(c.Request.PostForm)

		err := c.Request.ParseForm()
		if err != nil {

			c.String(500, fmt.Sprintf("ParseForm() err:%v", err))
			return
		}

		//c.String(200, "Post request successful\n")
		username := c.Request.FormValue("username")
		password := c.Request.FormValue("password")
		//c.String(200, "username=%s\n", username)
		//c.String(200, "password=%s\n", password)
		fmt.Printf("Username: %s\n", username)
		fmt.Printf("Password: %s\n", password)

		// Check if the username and password are valid
		if app.IsValidUser(db, username, password) {
			// If the username and password match, render the chatbot interaction page
			c.HTML(http.StatusOK, "./static/chatbot.html", gin.H{"username": username})
		} else {
			// If the username or password is incorrect, render the sign-in page again with an error message
			c.HTML(http.StatusOK, "./static/signin.html", gin.H{"error": "Incorrect username or password"})
		}
	})


chatbot.html

<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Chatbot Interaction</title>
    <!-- External CSS stylesheets -->
    <link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.5.2/css/bootstrap.min.css">
    <!-- Custom CSS styles -->
    <style>
        /* Body styles */
        body {
            font-family: Arial, sans-serif;
            background-color: #f8f9fa;
        }
        /* Container styles */
        .container {
            max-width: 800px;
            margin: 0 auto;
            padding: 20px;
        }
        /* Chat box styles */
        .chat-box {
            border: 1px solid #ced4da;
            padding: 20px;
            background-color: #fff;
            border-radius: 10px;
            box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
            min-height: 400px;
            overflow-y: auto;
        }
        /* Message styles */
        .message {
            margin-bottom: 15px;
            padding: 10px;
            border-radius: 10px;
        }
        /* User message styles */
        .user-message {
            background-color: #007bff;
            color: #fff;
            text-align: right;
        }
        /* Bot message styles */
        .bot-message {
            background-color: #28a745;
            color: #fff;
            text-align: left;
        }
        /* Input field and button styles */
        #messageInput {
            width: calc(100% - 80px);
            margin-bottom: 10px;
            padding: 10px;
            border-radius: 5px;
            border: 1px solid #ced4da;
            outline: none;
        }
        #sendBtn {
            width: 80px;
            padding: 10px;
            border: none;
            border-radius: 5px;
            background-color: #007bff;
            color: #fff;
            cursor: pointer;
            transition: background-color 0.3s;
        }
        #sendBtn:hover {
            background-color: #0056b3;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1 class="mb-4">Welcome, <span id="username">{{.username}}</span>!</h1>
        <div class="chat-box" id="chatBox">
            <!-- Chat messages will be displayed here -->
        </div>
        <form id="chatForm" class="mt-4">
            <input type="text" id="messageInput" placeholder="Type your message..." required>
            <button type="submit" id="sendBtn">Send</button>
        </form>
    </div>
    
    <!-- External JavaScript libraries -->
    <script src="https://code.jquery.com/jquery-3.5.1.slim.min.js"></script>
    <!-- Custom JavaScript code -->
    <script>
        // JavaScript code for handling chat interactions
        $(document).ready(function() {
            $("#chatForm").submit(function(event) {
                event.preventDefault(); // Prevent the default form submission

                // Get the user's message
                var message = $("#messageInput").val();

                // Display the user's message in the chat box
                displayMessage("user", message);

                // Process the user's message (You can send it to the server for processing)
                // For demonstration purposes, let's assume the bot's response is received here
                var botResponse = "Hello, <span style='font-weight: bold;'>{{.username}}</span>! How can I assist you?";
                
                // Display the bot's response in the chat box
                displayMessage("bot", botResponse);

                // Clear the input field
                $("#messageInput").val("");
            });
        });

        // Function to display messages in the chat box
        function displayMessage(sender, message) {
            var chatBox = $("#chatBox");
            var messageClass = sender === "user" ? "user-message" : "bot-message";
            var messageElement = $("<div>").addClass("message").addClass(messageClass).html(message);
            chatBox.append(messageElement);

            // Automatically scroll to the bottom of the chat box
            chatBox.scrollTop(chatBox[0].scrollHeight);
        }
    </script>
</body>
</html>


