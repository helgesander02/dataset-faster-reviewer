package services

import (
	"backend/src/utils"
	"log"
)

func (us *UserServices) SetBase64ImageCache(jobName string) {
	if !us.CacheManager.ExistsImageCacheStore(jobName) {
		us.CacheManager.SetImageCacheStore(jobName)
		log.Println("Image cache store created for job:", jobName)
	} else {
		log.Println("Image cache store already exists for job:", jobName)
	}
}

func (us UserServices) GetBase64ImageCacheByPage(jobName string, pageIndex int) ([]string, []string) {
	if !us.CacheManager.ExistsImageCacheStore(jobName) {
		us.CacheManager.SetImageCacheStore(jobName)
	}

	imageCache, exist := us.CacheManager.GetImageCacheStore(jobName)
	if !exist {
		log.Println("Image cache store not found for job:", jobName)
		return nil, nil
	}
	current_page_imagepath_set := us.CurrentPageData.GetPageItemAllImagePathByIndex(pageIndex)
	current_page_base64image_set := imageCache.GetBase64ImageCacheByImagePathSet(current_page_imagepath_set)

	return current_page_imagepath_set, current_page_base64image_set
}

func (us *UserServices) SetBase64ImageCacheByPage(jobName string, pageIndex int) ([]string, []string) {
	if !us.CacheManager.ExistsImageCacheStore(jobName) {
		us.CacheManager.SetImageCacheStore(jobName)
	}

	current_page_imagepath_set := us.CurrentPageData.GetPageItemAllImagePathByIndex(pageIndex)
	current_page_base64image_set := utils.CompressImageSetToBase64(current_page_imagepath_set)

	cacheData, found := us.CacheManager.GetImageCacheStore(jobName)
	if !found {
		log.Println("Image cache store not found for job:", jobName)
		return nil, nil
	}

	cacheData.SetBase64ImageCacheByImagePathSet(current_page_imagepath_set, current_page_base64image_set)
	us.CacheManager.UpdateImageCacheStore(jobName, cacheData)
	return current_page_imagepath_set, current_page_base64image_set
}

func (us *UserServices) ImageCacheExists(jobName string) bool {
	_, found := us.CacheManager.GetImageCacheStore(jobName)
	if !found {
		log.Println("Image cache store not found for job:", jobName)
		return false
	}
	return true
}
