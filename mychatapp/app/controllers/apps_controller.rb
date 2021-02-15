class AppsController < ApplicationController
  def index
      limit = Integer(params[:limit]) || 10
      offset = ((Integer(params[:page]) || 1) - 1) * limit
      apps = App.order('created_at DESC').limit(limit).offset(offset);
      render json: {status: 'SUCCESS', message:'Loaded apps', data:apps},status: :ok
  end

  def show
      app = App.find(params[:token])
      render json: {status: 'SUCCESS', message:'Loaded app', data:app},status: :ok
  end

  def create
      app = App.new(app_params)
      if app.save
        Redis.current.set(app.token, app.id)
        render json: {status: 'SUCCESS', message:'Saved app', data:app},status: :ok
      else
        render json: {status: 'ERROR', message:'Failed to create app', data:app.errors},status: :unprocessable_entity
      end
  end

  def update
      app = App.find(params[:id])
      if app.update_attributes(app_params)
        render json: {status: 'SUCCESS', message:'Updated app', data:app},status: :ok
      else
        render json: {status: 'ERROR', message:'Failed to update app', data:app.errors},status: :unprocessable_entity
      end
  end

  def app_params
      params.require(:app).permit(:name)
  end
end

