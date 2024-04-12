use std::env;
use chrono::Utc;
use tonic::{transport::Server, Request, Response, Status};
pub mod grpc_chat {
    tonic::include_proto!("chat");
}
use grpc_chat::chat_service_server::{ChatService, ChatServiceServer};
use grpc_chat::{ChatMessage, UserData};
#[derive(Debug, Default)]
pub struct ChatServiceRPC {}
#[tonic::async_trait]
impl ChatService for ChatServiceRPC {
    type SendMessageStream = tonic::Streaming<ChatMessage>;
    type ReceiveMessagesStream = tonic::Streaming<ChatMessage>;

    async fn send_message(
        &self,
        request: Request<ChatMessage>,
    ) -> Result<Response<Self::SendMessageStream>, Status> {
        println!("Got a request: {:?}", request);
        let mut chats = Vec::new();        
        chats.push(grpc_chat::ChatMessage{
            sender_id: 1.to_string(),
            reciever_id: 2.to_string(),
            content: "hello".to_string(),
            timestamp: timestamp,
        });
    }

    async fn receive_messages(
        &self,
        request: Request<UserData>,
    ) -> Result<Response<Self::ReceiveMessagesStream>, Status> {
        println!("Got a request: {:?}", request);
        // Implement receive_messages logic here
        unimplemented!()
    }
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let port = env::var("PORT").unwrap_or("50051".to_string());
    let addr = format!("0.0.0.0:{}", port).parse()?;
    let chat_service = ChatServiceRPC::default();
    let chat_service = ChatServiceServer::new(chat_service);
    let chat_service = tonic_web::enable(chat_service);
    let (mut health_reporter, health_service) = tonic_health::server::health_reporter();
    health_reporter
        .set_serving::<ChatServiceServer<ChatServiceRPC>>()
        .await;
    println!("Running on port {}...", port);
    Server::builder()
        .accept_http1(true)
        .add_service(health_service)
        .add_service(chat_service)
        .serve(addr)
        .await?;
    Ok(())
}
