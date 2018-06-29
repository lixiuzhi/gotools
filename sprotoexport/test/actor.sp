
// 移除角色原因
enum KillActorReason{
	None  	 		
	AreaDamage 		
	BulletDamage 	
	Respawn  		
	SkillEnd 		
	InstanceStop	
	TriggerActionVanish
	LifeTimeOut
	Buff
	Collect // 采集
	GMReSpawn 
	GMKill	// GM指令杀死目标
	
}


//{"Uid":100}
message ActorSyncPosTypeACK{
	
	ObjID int64
	
}



//角色坐骑变化
//{"Uid":200}
message ActorCarACK{
	
	ObjID int64
	
	CarID int32		// 坐骑ID 0:下马 >0:上马
	
}

// 角色装扮变化
//{"Uid":2020}
message ActorAvatarACK{
	
	ObjID int64
	AvatarItemID	[]int32			// 装扮列表
}

