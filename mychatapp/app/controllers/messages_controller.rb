class MessagesController < ApplicationController

  def index
      token = params[:token]
      begin
        appId = Redis.current.get(token)
      rescue 
        render json: {status: 'Failed', message:'token not foumd', data:messages},status: :unprocessable_entity
      
      chat = Chat.where("apps_id = :appId and chat_number = :chatNumber" ,{ appId: appId, chatNumber: chat_number })
      if !chat 
        render json: {status: 'Failed', message:'chat not found', data:messages},status: :unprocessable_entity
      
      limit = Integer(params[:limit]) || 10
      offset = ((Integer(params[:page]) || 1) - 1) * limit
      messages = Message.where("chats_id = :chatId" ,{ chatId: chat.id }).order('message_number DESC').limit(limit).offset(offset);
      render json: {status: 'SUCCESS', message:'Loaded chats', data:messages},status: :ok
  end

end
