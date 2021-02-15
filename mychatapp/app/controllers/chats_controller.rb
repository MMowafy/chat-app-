class ChatsController < ApplicationController
  def index
        token = params[:token]
        appId = Redis.current.get(token)
        limit = Integer(params[:limit]) || 10
        offset = ((Integer(params[:page]) || 1) - 1) * limit
        chats = Chat.where("apps_id = :appId" ,{ appId: appId }).order('chat_number DESC').limit(limit).offset(offset);
        render json: {status: 'SUCCESS', message:'Loaded chats', data:chats},status: :ok
  end

end
